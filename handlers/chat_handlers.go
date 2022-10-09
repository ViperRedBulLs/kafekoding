package handlers

import (
	"kafekoding/models"
	"kafekoding/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ChatIndex(ctx *gin.Context) {
	session := sessions.Default(ctx)

	user := session.Get("user")
	if user == nil {
		flash.Set("danger", "Please login first before accessing this page.")
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	if ctx.Request.Method == http.MethodGet {
		rooms := models.NewChatModel().GetAllRooms()

		username := utils.GetSessionValue(user)["username"].(string)
		var chatAction ChatAction
		err := ctx.Bind(&chatAction)
		if err == nil {
			if chatAction.Type == "join" {
				models.NewChatModel().AddMember(chatAction.RoomID, username)
				flash.Set("success", "Successfully to join room.")
				ctx.Redirect(http.StatusFound, "/chat")
				return
			}
			log.Println(chatAction)
		}

		defer flash.Delete()

		ctx.HTML(http.StatusOK, "chat_index", gin.H{
			"user":  user,
			"rooms": rooms,
			"flash": flash,
		})
		return
	}
}

// ChatDetail
func ChatDetail(ctx *gin.Context) {
	session := sessions.Default(ctx)

	user := session.Get("user")
	if user == nil {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	if ctx.Request.Method == http.MethodGet {
		roomID := ctx.Param("roomID")
		roomIDInt, _ := strconv.Atoi(roomID)

		room, err := models.NewChatModel().GetRoomByID(uint(roomIDInt))
		if err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusNotFound)
			return
		}

		// get username and return by json
		if username, ok := ctx.GetQuery("check-user-fullname"); ok {
			u, err := models.NewUserModel().GetUserByUsername(username)
			if err != nil {
				ctx.JSON(http.StatusNotFound, nil)
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": u.FirstName + " " + u.LastName,
			})
			return
		}

		chats := models.NewChatModel().GetAllChatsByRoomID(room.ID)
		userAdmin, _ := models.NewUserModel().GetUserByID(room.AdminID)

		ctx.HTML(http.StatusOK, "chat_detail", gin.H{
			"user":      user,
			"room":      room,
			"chats":     chats,
			"userAdmin": userAdmin,
		})
		return
	}

	if ctx.Request.Method == http.MethodPost {
		roomID := ctx.PostForm("room_id")
		senderID := ctx.PostForm("sender_id")
		text := ctx.PostForm("text")

		roomIDInt, _ := strconv.Atoi(roomID)
		senderIDInt, _ := strconv.Atoi(senderID)

		chat := models.ChatMessage{
			RoomID:   uint(roomIDInt),
			SenderID: uint(senderIDInt),
			Text:     text,
		}

		err := models.NewChatModel().CreateChat(&chat)
		if err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Successfully to send message.",
		})
		return
	}
}
