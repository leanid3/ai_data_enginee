.PHONY: proto build run clean test

# Генерация proto файлов
proto:
	protoc --proto_path=. --proto_path=third_party \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		proto/file.proto

# Сборка
build:
	go build -o main ./cmd/main.go

# Запуск
run: build
	./main

# Очистка
clean:
	rm -f main

# Тесты
test:
	go test ./...

# Установка зависимостей
deps:
	go mod download
	go mod tidy