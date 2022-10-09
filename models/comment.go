package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CourseID  uint           `json:"course_id"`
	UserID    uint           `json:"user_id"`
	File      string         `json:"file" gorm:"size:255;null"`
	Content   string         `gorm:"type:longtext;null" json:"content"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Replies   []Reply        `gorm:"foreignKey:CommentID" json:"replies"`
}

type CommentModel struct {
	db *gorm.DB
}

func NewCommentModel() *CommentModel {
	return &CommentModel{getDB()}
}

func (c *CommentModel) GetAllComments() []Comment {
	var comments []Comment
	c.db.Model(&Comment{}).Preload("Replies").Find(&comments)
	return comments
}
