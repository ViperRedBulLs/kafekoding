package models

import (
	"kafekoding/utils"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                 uint        `gorm:"primaryKey" json:"id"`
	FirstName          string      `gorm:"size:100" json:"first_name"`
	LastName           string      `gorm:"size:100" json:"last_name"`
	Username           string      `gorm:"size:100;unique;index" json:"username"`
	Email              string      `gorm:"size:100;unique;index" json:"email"`
	Password           string      `gorm:"size:255"`
	IsSuperuser        bool        `gorm:"default:0" json:"is_superuser"`
	IsStaff            bool        `gorm:"default:0" json:"is_staff"`
	IsActive           bool        `gorm:"default:1" json:"is_active"`
	IsOnline           bool        `gorm:"default:0" json:"is_online"`
	LastLogin          time.Time   `json:"last_login" gorm:"null"`
	DateJoined         time.Time   `gorm:"autoCreateTime" json:"date_joined"`
	RoomMentors        []Room      `gorm:"foreignKey:MentorID" json:"room_mentor"`
	RoomMembers        []*Room     `gorm:"many2many:user_room_members" json:"room_members"`
	ChatRooms          []*ChatRoom `gorm:"many2many:user_chatrooms_members" json:"chat_rooms"`
	CourseUserComments []Comment   `gorm:"foreignKey:UserID" json:"course_user_comment"`
}

type UserModel struct {
	db *gorm.DB
}

func NewUserModel() *UserModel {
	return &UserModel{getDB()}
}

func (u *UserModel) CreateUser(user *User) error {
	user.Password = utils.EncryptionPassword(user.Password)
	err := u.db.Create(user).Error
	return err
}

func (u *UserModel) GetUserByID(id uint) (User, error) {
	var user User
	err := u.db.Model(&User{}).Where("id = ?", id).Preload("RoomMentors").Preload("RoomMembers").Preload("CourseUserComments").First(&user).Error
	return user, err
}

func (u *UserModel) GetUserByUsername(username string) (User, error) {
	var user User
	err := u.db.Model(&User{}).Where("username = ?", username).Preload("RoomMentors").Preload("RoomMembers").Preload("CourseUserComments").First(&user).Error
	return user, err
}

func (u *UserModel) GetAllUsers() []User {
	var users []User
	u.db.Model(&User{}).Preload("RoomMentors").Preload("RoomMembers").Preload("CourseUserComments").Find(&users)
	return users
}
