package repository

import (
	"context"

	"github.com/AlGrushino/chat/internal/repository/chat"
	"github.com/AlGrushino/chat/internal/repository/message"
	"github.com/AlGrushino/chat/internal/repository/models"
	"gorm.io/gorm"
)

type Message interface {
	Create(ctx context.Context, message *models.Message) error
	GetByChatID(ctx context.Context, chatID int, limit, offset int) ([]*models.Message, error)
	GetByID(ctx context.Context, id int) (*models.Message, error)
	Delete(ctx context.Context, id int) error
}

type Chat interface {
	Create(ctx context.Context, chat *models.Chat) error
	GetByID(ctx context.Context, id int) (*models.Chat, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Chat, error)
	Update(ctx context.Context, chat *models.Chat) error
	Delete(ctx context.Context, id int) error
	ChatExists(ctx context.Context, title string) (bool, error)
}

type Repository struct {
	Chat
	Message
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Chat:    chat.NewChatRepository(db),
		Message: message.NewMessageRepository(db),
	}
}
