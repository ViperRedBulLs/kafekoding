package models

import (
	"time"

	"gorm.io/gorm"
)

type Reply struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CommentID uint           `json:"comment_id"`
	UserID    uint           `json:"user_id"`
	File      string         `gorm:"size:255;null" json:"file"`
	Content   string         `gorm:"type:longtext;null" json:"content"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
