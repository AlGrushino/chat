package service

import (
	"chat/internal/repository"
	"chat/internal/service/chat"
	"context"

	"github.com/sirupsen/logrus"
)

type Chat interface {
	CreateChat(ctx context.Context, title string) (string, error)
}

type Service struct {
	Chat
}

func NewService(log *logrus.Logger, repository *repository.Repository) *Service {
	return &Service{
		Chat: chat.NewChatService(log, repository),
	}
}
