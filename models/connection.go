package models

import (
	"log"

	gormsession "github.com/gin-contrib/sessions/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GetDB = getDB()

func init() {
	db := getDB()

	err := db.AutoMigrate(&User{}, &Room{}, &Course{}, &Comment{}, &Reply{}, &ChatRoom{},
		&ChatMessage{})
	if err != nil {
		return
	}
}

func getDB() *gorm.DB {
	dsn := "root:@/kafekoding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetStore() gormsession.Store {
	store := gormsession.NewStore(getDB(), true, []byte("secret"))
	return store
}
