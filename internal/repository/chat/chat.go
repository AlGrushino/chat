package chat

import (
	"context"
	"fmt"

	"github.com/AlGrushino/chat/internal/repository/models"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(ctx context.Context, chat *models.Chat) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

func (r *ChatRepository) GetByID(ctx context.Context, id int) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.WithContext(ctx).First(&chat, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get chat by id: %w", err)
	}

	return &chat, nil
}

func (r *ChatRepository) ChatExists(ctx context.Context, title string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Chat{}).
		Where("title = ?", title).
		Count(&count).Error

	return count > 0, err
}

func (r *ChatRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Chat, error) {
	var chats []*models.Chat
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&chats).Error
	return chats, err
}

func (r *ChatRepository) Update(ctx context.Context, chat *models.Chat) error {
	return r.db.WithContext(ctx).Save(chat).Error
}

func (r *ChatRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.Chat{}, id)
	return result.Error
}
