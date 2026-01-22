package message

import (
	"context"
	"errors"
	"fmt"
	"sort"
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

func (s *MessageService) GetMessages(ctx context.Context, id, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		return nil, fmt.Errorf("limit is too big: %d", limit)
	}

	_, err := s.repository.Chat.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chat does not exist")
		}
		return nil, fmt.Errorf("failed to get chat with id: %d", id)
	}

	messages, err := s.repository.Message.GetByChatID(ctx, id, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	sort.Slice(messages, func(i, j int) bool { return messages[i].CreatedAt.Before(messages[j].CreatedAt) })

	messageList := make([]string, 0, len(messages))
	for _, message := range messages {
		messageList = append(messageList, message.Text)
	}

	return messageList, nil
}
