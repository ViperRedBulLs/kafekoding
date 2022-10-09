package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RoomID    uint           `json:"room_id"`
	Title     string         `gorm:"size:100;unique;index" json:"title"`
	Content   string         `gorm:"type:longtext;null" json:"content"`
	File      string         `gorm:"size:255;null" json:"file"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Comments  []Comment      `gorm:"foreignKey:CourseID" json:"comments"`
}

type CourseModel struct {
	db *gorm.DB
}

func NewCourseModel() *CourseModel {
	return &CourseModel{getDB()}
}

func (c *CourseModel) CreateCourse(course *Course) error {
	err := c.db.Create(course).Error
	return err
}

func (c *CourseModel) GetCourseByID(id uint) (Course, error) {
	var course Course
	err := c.db.Model(&Course{}).Where("id = ?", id).Preload("Comments").First(&course).Error
	return course, err
}

func (c *CourseModel) GetCourseByTitle(title string) (Course, error) {
	var course Course
	err := c.db.Model(&Course{}).Where("title = ?", title).Preload("Comments").First(&course).Error
	return course, err
}

func (c *CourseModel) GetAllCourses() []Course {
	var courses []Course
	c.db.Model(&Course{}).Preload("Comments").Find(&courses)
	return courses
}
