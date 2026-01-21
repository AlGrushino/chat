package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	ID        int       `gorm:"primaryKey"`
	Title     string    `gorm:"size:200;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Message struct {
	ID        int       `gorm:"primaryKey"`
	ChatID    int       `gorm:"not null"`
	Text      string    `gorm:"size:5000;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Chat Chat `gorm:"foreignKey:ChatID"`
}

func (c *Chat) BeforeCreate(tx *gorm.DB) (err error) {
	if len(c.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters")
	}
	return nil
}

func (c *Chat) BeforeUpdate(tx *gorm.DB) (err error) {
	if len(c.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters")
	}
	return nil
}
