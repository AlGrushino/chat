package service

import (
	"context"

	"github.com/AlGrushino/chat/internal/repository"
	"github.com/AlGrushino/chat/internal/service/chat"
	"github.com/AlGrushino/chat/internal/service/message"
	"github.com/sirupsen/logrus"
)

type Chat interface {
	CreateChat(ctx context.Context, title string) (string, error)
	DeleteChat(ctx context.Context, id int) error
}

type Message interface {
	AddMessage(ctx context.Context, id int, text string) (string, error)
	GetMessages(ctx context.Context, id, limit int) ([]string, error)
}

type Service struct {
	Chat
	Message
}

func NewService(log *logrus.Logger, repository *repository.Repository) *Service {
	return &Service{
		Chat:    chat.NewChatService(log, repository),
		Message: message.NewMessageService(log, repository),
	}
}
