package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;unique;index" json:"name"`
	MentorID  uint           `json:"mentor_id"`
	Members   []*User        `gorm:"many2many:user_room_members;" json:"members"`
	Logo      string         `gorm:"size:255;null" json:"logo"`
	Desc      string         `gorm:"size:255;null" json:"desc"`
	Content   string         `gorm:"type:longtext;null" json:"content"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Courses   []Course       `gorm:"foreignKey:RoomID" json:"courses"`
}

type RoomModel struct {
	db *gorm.DB
}

func NewRoomModel() *RoomModel {
	return &RoomModel{getDB()}
}

func (r *RoomModel) CreateRoom(room *Room) error {
	err := r.db.Create(room).Error
	return err
}

func (r *RoomModel) GetRoomByID(id uint) (Room, error) {
	var room Room
	err := r.db.Model(&Room{}).Where("id = ?", id).Preload("Members").First(&room).Error
	return room, err
}

func (r *RoomModel) GetRoomByName(name string) (Room, error) {
	var room Room
	err := r.db.Model(&Room{}).Where("name = ?", name).Preload("Members").First(&room).Error
	return room, err
}

func (r *RoomModel) GetAllRooms() []Room {
	var rooms []Room
	r.db.Model(&Room{}).Preload("Members").Find(&rooms)
	return rooms
}

func (r *RoomModel) AddMember(roomID uint, username string) {

	user, _ := NewUserModel().GetUserByUsername(username)
	room, _ := r.GetRoomByID(roomID)

	room.Members = append(room.Members, &user)
	r.db.Save(&room)
}

func (r *RoomModel) RemoveMember(roomID uint, username string) {
	user, _ := NewUserModel().GetUserByUsername(username)
	room, _ := r.GetRoomByID(roomID)

	r.db.Model(&room).Association("Members").Delete(&user)
}
