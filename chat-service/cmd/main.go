package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Простые структуры для тестирования
type CreateDialogRequest struct {
	UserId         string `json:"user_id"`
	Title          string `json:"title"`
	InitialMessage string `json:"initial_message"`
}

type CreateDialogResponse struct {
	DialogId  string `json:"dialog_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type SendMessageRequest struct {
	DialogId string `json:"dialog_id"`
	UserId   string `json:"user_id"`
	Message  string `json:"message"`
}

type SendMessageResponse struct {
	MessageId          string `json:"message_id"`
	Status             string `json:"status"`
	Message            string `json:"message"`
	CreatedAt          string `json:"created_at"`
	RequiresProcessing bool   `json:"requires_processing"`
}

type GetDialogHistoryRequest struct {
	DialogId string `json:"dialog_id"`
	UserId   string `json:"user_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type GetDialogHistoryResponse struct {
	DialogId   string     `json:"dialog_id"`
	Messages   []*Message `json:"messages"`
	TotalCount int32      `json:"total_count"`
	Page       int32      `json:"page"`
	PageSize   int32      `json:"page_size"`
}

type ListDialogsRequest struct {
	UserId   string `json:"user_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type ListDialogsResponse struct {
	Dialogs    []*DialogInfo `json:"dialogs"`
	TotalCount int32         `json:"total_count"`
	Page       int32         `json:"page"`
	PageSize   int32         `json:"page_size"`
}

type DeleteDialogRequest struct {
	DialogId string `json:"dialog_id"`
	UserId   string `json:"user_id"`
}

type DeleteDialogResponse struct {
	DialogId string `json:"dialog_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type UpdateMessageStatusRequest struct {
	MessageId string `json:"message_id"`
	Status    string `json:"status"`
	Result    string `json:"result"`
}

type UpdateMessageStatusResponse struct {
	MessageId string `json:"message_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type GetMessageStatusRequest struct {
	MessageId string `json:"message_id"`
}

type GetMessageStatusResponse struct {
	MessageId string `json:"message_id"`
	Status    string `json:"status"`
	Result    string `json:"result"`
	UpdatedAt string `json:"updated_at"`
}

type Message struct {
	MessageId string `json:"message_id"`
	Content   string `json:"content"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type DialogInfo struct {
	DialogId     string `json:"dialog_id"`
	Title        string `json:"title"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	MessageCount int32  `json:"message_count"`
}

func main() {
	// Создаем HTTP сервер
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/v1/dialogs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Создание диалога
			var req CreateDialogRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response := CreateDialogResponse{
				DialogId:  uuid.New().String(),
				Status:    "created",
				Message:   "Диалог создан",
				CreatedAt: time.Now().Format(time.RFC3339),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/v1/dialogs/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Получение истории диалога
			response := GetDialogHistoryResponse{
				DialogId:   "test-dialog-id",
				Messages:   []*Message{},
				TotalCount: 0,
				Page:       1,
				PageSize:   10,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Chat Service запущен на порту 50055")
	log.Fatal(http.ListenAndServe(":50055", nil))
}
