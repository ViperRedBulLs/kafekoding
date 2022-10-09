package apiv1

import (
	"fmt"
	"kafekoding/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type APIV1Room interface {
	GetAllRooms() gin.HandlerFunc
	GetRoomByName() gin.HandlerFunc
}

type apiV1Room struct{}

func NewAPIV1Room() APIV1Room {
	return &apiV1Room{}
}

func (a apiV1Room) GetAllRooms() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user := session.Get("user")
		if user == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  "error",
				"message": "Please login first, before accessing this page.",
			})
			return
		}

		roomAPI := []map[string]interface{}{}
		rooms := models.NewRoomModel().GetAllRooms()
		for _, room := range rooms {
			members := []map[string]interface{}{}
			for _, m := range room.Members {
				members = append(members, map[string]interface{}{
					"id":       m.ID,
					"username": m.Username,
				})
			}
			roomAPI = append(roomAPI, map[string]interface{}{
				"id":         room.ID,
				"name":       room.Name,
				"mentor_id":  room.MentorID,
				"members":    members,
				"desc":       room.Desc,
				"content":    room.Content,
				"created_at": room.CreatedAt,
				"updated_at": room.UpdatedAt,
				"deleted_at": room.DeletedAt,
			})
		}
		ctx.JSON(http.StatusOK, roomAPI)
	}
}

func (a apiV1Room) GetRoomByName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user := session.Get("user")
		if user == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Please login firrst, before accessing this page.",
			})
			return
		}

		name := ctx.Param("name")
		room, err := models.NewRoomModel().GetRoomByName(name)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": fmt.Sprintf("Error! Room by %s not found.", name),
			})
			return
		}

		members := []map[string]interface{}{}
		for _, member := range room.Members {
			members = append(members, map[string]interface{}{
				"id":       member.ID,
				"username": member.Username,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"id":         room.ID,
			"name":       room.Name,
			"mentor_id":  room.MentorID,
			"members":    members,
			"desc":       room.Desc,
			"content":    room.Content,
			"updated_at": room.UpdatedAt,
			"created_at": room.CreatedAt,
			"deleted_at": room.DeletedAt,
		})
	}
}
