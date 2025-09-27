package main

import (
	"context"
	"log"
	"net"
	"net/http"

	chatGen "backend/api-gateway/chat-gen"
	fileGen "backend/api-gateway/file-gen"
	"backend/api-gateway/gen"
	"backend/api-gateway/internal/handlers"
	llmGen "backend/api-gateway/llm-gen"
	orchestrationGen "backend/api-gateway/orchestration-gen"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Создаем TCP listener для gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Подключение к новым микросервисам
	var fileConn *grpc.ClientConn
	var chatConn *grpc.ClientConn
	var llmConn *grpc.ClientConn
	var orchestrationConn *grpc.ClientConn

	// Подключаемся к File Service
	fileConn, err = grpc.Dial("file-service:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*1024), // 1GB
			grpc.MaxCallSendMsgSize(1024*1024*1024), // 1GB
		),
	)
	if err != nil {
		log.Fatalf("failed to dial file service: %v", err)
	}
	fileClient := fileGen.NewFileServiceClient(fileConn)

	// Подключаемся к Chat Service
	chatConn, err = grpc.Dial("chat-service:50055",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial chat service: %v", err)
	}
	chatClient := chatGen.NewChatServiceClient(chatConn)

	// Подключаемся к LLM Service
	llmConn, err = grpc.Dial("llm-service:50056",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial llm service: %v", err)
	}
	llmClient := llmGen.NewLLMServiceClient(llmConn)

	// Подключаемся к Orchestration Service
	orchestrationConn, err = grpc.Dial("orchestration-service:50057",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial orchestration service: %v", err)
	}
	orchestrationClient := orchestrationGen.NewOrchestrationServiceClient(orchestrationConn)

	// Настраиваем gRPC сервер для больших сообщений
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024*1024), // 1GB
		grpc.MaxSendMsgSize(1024*1024*1024), // 1GB
	)

	// Создаем обработчик с новыми клиентами
	gatewayHandler := handlers.NewGatewayHandler(fileClient, chatClient, llmClient, orchestrationClient)
	gen.RegisterDataEngineeringAssistantServer(grpcServer, gatewayHandler)

	go func() {
		log.Println("Started gRPC server :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = gen.RegisterDataEngineeringAssistantHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	// Создаем основной HTTP роутер
	httpMux := http.NewServeMux()

	// Добавляем gRPC Gateway роуты
	httpMux.Handle("/", mux)

	// Добавляем health check
	httpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Добавляем маршруты для анализа данных
	analysisHandler := handlers.NewAnalysisHandler("http://airflow-webserver:8080", "http://ollama:11434")
	analysisHandler.RegisterAnalysisRoutes(httpMux)

	log.Println("Started HTTP server :8080")
	http.ListenAndServe(":8080", httpMux)
}
