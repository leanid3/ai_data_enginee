package handlers

import (
	"context"
	"time"

	"chat-service/gen"

	"github.com/google/uuid"
)

type ChatHandler struct {
	gen.UnimplementedChatServiceServer
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

func (h *ChatHandler) CreateDialog(ctx context.Context, req *gen.CreateDialogRequest) (*gen.CreateDialogResponse, error) {
	dialogID := uuid.New().String()

	return &gen.CreateDialogResponse{
		DialogId:  dialogID,
		Status:    "created",
		Message:   "Диалог успешно создан",
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func (h *ChatHandler) SendMessage(ctx context.Context, req *gen.SendMessageRequest) (*gen.SendMessageResponse, error) {
	messageID := uuid.New().String()

	return &gen.SendMessageResponse{
		MessageId:          messageID,
		Status:             "sent",
		Message:            "Сообщение отправлено",
		CreatedAt:          time.Now().Format(time.RFC3339),
		RequiresProcessing: true,
	}, nil
}

func (h *ChatHandler) GetDialogHistory(ctx context.Context, req *gen.GetDialogHistoryRequest) (*gen.GetDialogHistoryResponse, error) {
	return &gen.GetDialogHistoryResponse{
		DialogId:   req.DialogId,
		Messages:   []*gen.Message{},
		TotalCount: 0,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

func (h *ChatHandler) ListDialogs(ctx context.Context, req *gen.ListDialogsRequest) (*gen.ListDialogsResponse, error) {
	return &gen.ListDialogsResponse{
		Dialogs:    []*gen.DialogInfo{},
		TotalCount: 0,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

func (h *ChatHandler) DeleteDialog(ctx context.Context, req *gen.DeleteDialogRequest) (*gen.DeleteDialogResponse, error) {
	return &gen.DeleteDialogResponse{
		DialogId: req.DialogId,
		Status:   "deleted",
		Message:  "Диалог удален",
	}, nil
}

func (h *ChatHandler) UpdateMessageStatus(ctx context.Context, req *gen.UpdateMessageStatusRequest) (*gen.UpdateMessageStatusResponse, error) {
	return &gen.UpdateMessageStatusResponse{
		MessageId: req.MessageId,
		Status:    "updated",
		Message:   "Статус обновлен",
	}, nil
}

func (h *ChatHandler) GetMessageStatus(ctx context.Context, req *gen.GetMessageStatusRequest) (*gen.GetMessageStatusResponse, error) {
	return &gen.GetMessageStatusResponse{
		MessageId: req.MessageId,
		Status:    "completed",
		Result:    "Обработка завершена",
		UpdatedAt: time.Now().Format(time.RFC3339),
	}, nil
}
