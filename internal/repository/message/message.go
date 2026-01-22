package message

import (
	"context"

	"github.com/AlGrushino/chat/internal/repository/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *MessageRepository) GetByChatID(ctx context.Context, chatID int, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message
	err := r.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Limit(limit).
		Offset(offset).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) GetByID(ctx context.Context, id int) (*models.Message, error) {
	var message models.Message
	err := r.db.WithContext(ctx).
		Preload("Chat").
		First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Message{}, id).Error
}
