package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type ChatRoom struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	AdminID uint    `json:"admin_id"`
	Members []*User `gorm:"many2many:user_chatrooms_members" json:"members"`
	Name    string  `gorm:"size:100;unique;index"`
}

type ChatMessage struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	RoomID   uint      `json:"room_id"`
	SenderID uint      `json:"sender_id"`
	File     string    `json:"file" gorm:"size:255;null"`
	Text     string    `gorm:"type:longtext;null" json:"text"`
	SendAt   time.Time `gorm:"autoCreateTime" json:"send_at"`
}

type ChatModel interface {
	CreateRoom(room *ChatRoom) error
	GetAllRooms() []ChatRoom
	GetRoomByName(name string) (ChatRoom, error)
	GetRoomByID(id uint) (ChatRoom, error)
	AddMember(roomID uint, username string)
	CreateChat(chat *ChatMessage) error
	GetAllChatsByRoomID(roomID uint) []ChatMessage
}

type chatModel struct {
	db *gorm.DB
}

func NewChatModel() ChatModel {
	return &chatModel{getDB()}
}

func (c *chatModel) CreateRoom(room *ChatRoom) error {
	err := c.db.Create(room).Error
	return err
}

func (c *chatModel) GetAllRooms() []ChatRoom {
	var rooms []ChatRoom
	c.db.Model(&ChatRoom{}).Preload("Members").Find(&rooms)
	return rooms
}

func (c *chatModel) GetRoomByName(name string) (ChatRoom, error) {
	var room ChatRoom
	err := c.db.Model(&ChatRoom{}).Where("name = ?", name).Preload("Members").First(&room).Error
	return room, err
}

func (c *chatModel) GetRoomByID(id uint) (ChatRoom, error) {
	var room ChatRoom
	err := c.db.Model(&ChatRoom{}).Where("id = ?", id).Preload("Members").First(&room).Error
	return room, err
}

func (c *chatModel) AddMember(roomID uint, username string) {
	user, err := NewUserModel().GetUserByUsername(username)
	if err != nil {
		log.Println(err.Error())
		return
	}
	room, err := c.GetRoomByID(roomID)
	if err != nil {
		log.Println(err.Error())
		return
	}

	room.Members = append(room.Members, &user)

	c.db.Save(&room)
}

func (c *chatModel) CreateChat(chat *ChatMessage) error {
	err := c.db.Create(chat).Error
	return err
}

func (c *chatModel) GetAllChatsByRoomID(roomID uint) []ChatMessage {
	var chats []ChatMessage
	c.db.Model(&ChatMessage{}).Where("room_id", roomID).Find(&chats)
	return chats
}
