package chat

import (
	"chat/internal/repository"
	"chat/internal/repository/models"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ChatService struct {
	repository *repository.Repository
	log        *logrus.Logger
}

func NewChatService(log *logrus.Logger, repository *repository.Repository) *ChatService {
	return &ChatService{
		repository: repository,
		log:        log,
	}
}

func (s *ChatService) CreateChat(ctx context.Context, title string) (string, error) {
	trimmed := strings.TrimSpace(title)

	length := len(trimmed)
	if length < 1 {
		s.log.Warnf("Create chat failed: empty title (original: %q)", title)
		return "", fmt.Errorf("len of %s equals 0", trimmed)
	}

	if length > 200 {
		s.log.Warnf("Create chat failed: title too long %d chars", len(trimmed))
		return "", fmt.Errorf("len of %s greater than 200", trimmed)
	}

	chat := models.Chat{
		Title:     trimmed,
		CreatedAt: time.Now(),
	}

	exist, err := s.repository.Chat.ChatExists(ctx, trimmed)
	if err != nil {
		s.log.WithError(err).Error("Failed to check if chat exists in database")
		return "", fmt.Errorf("failed to check if chat exists: %w", err)
	}

	if exist {
		s.log.Warnf("Create chat failed: title: %s already exists", trimmed)
		return "", fmt.Errorf("failed to create chat, chat: %s already exists", trimmed)
	}

	if err := s.repository.Chat.Create(ctx, &chat); err != nil {
		s.log.WithError(err).Error("Failed to create chat in database")
		return "", fmt.Errorf("failed to create chat: %w", err)
	}

	s.log.Infof("Chat created successfully (ID: %d, Title: %q)", chat.ID, trimmed)
	return trimmed, nil
}
