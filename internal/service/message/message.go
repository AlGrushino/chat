package message

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlGrushino/chat/internal/repository"
	"github.com/AlGrushino/chat/internal/repository/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MessageService struct {
	repository *repository.Repository
	log        *logrus.Logger
}

func NewMessageService(log *logrus.Logger, repository *repository.Repository) *MessageService {
	return &MessageService{
		repository: repository,
		log:        log,
	}
}

func (s *MessageService) AddMessage(ctx context.Context, id int, text string) (string, error) {
	_, err := s.repository.Chat.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("chat does not exist")
		}
		return "", fmt.Errorf("failed to get chat with id: %d", id)
	}

	if text == "" {
		return "", errors.New("text of message is empty")
	}

	if len(text) > 5000 {
		return "", errors.New("text of message is too long")
	}

	message := models.Message{
		ChatID:    id,
		Text:      text,
		CreatedAt: time.Now(),
	}

	err = s.repository.Message.Create(ctx, &message)
	if err != nil {
		return "", fmt.Errorf("failed to add message: %w", err)
	}

	return text, nil
}
