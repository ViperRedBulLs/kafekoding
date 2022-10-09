package handlers

import (
	"kafekoding/models"
	"kafekoding/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RoomIndex() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user := session.Get("user")
		if user == nil {
			flash.Set("info", "Please login first before accessing this page.")
			ctx.Redirect(http.StatusFound, "/login")
			return
		}

		// Method GET
		if ctx.Request.Method == http.MethodGet {
			defer flash.Delete()

			rooms := models.NewRoomModel().GetAllRooms()
			ctx.HTML(http.StatusOK, "index", gin.H{
				"flash": flash,
				"user":  user,
				"rooms": rooms,
			})
			return
		}

		// Method POST
		if ctx.Request.Method == http.MethodPost {
			mentor_id := utils.GetSessionValue(user)["id"].(uint)
			name := ctx.PostForm("name")
			logo, err := ctx.FormFile("logo")
			if err != nil {
				http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
				return
			}
			desc := ctx.PostForm("desc")
			content := ctx.PostForm("content")

			filename := logo.Filename

			room := models.Room{
				MentorID: mentor_id,
				Name:     name,
				Logo:     "/media/rooms/" + filename,
				Desc:     desc,
				Content:  content,
			}

			err = models.NewRoomModel().CreateRoom(&room)
			if err != nil {
				flash.Set("danger", err.Error())
				ctx.Redirect(http.StatusFound, "/")
				return
			}

			err = ctx.SaveUploadedFile(logo, "media/rooms/"+filename)
			if err != nil {
				http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
				return
			}

			flash.Set("success", "Successfully create new room.")
			ctx.Redirect(http.StatusFound, "/")
			return
		}
	}
}

func RoomDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user := session.Get("user")
		if user == nil {
			ctx.Redirect(http.StatusFound, "/login")
			return
		}

		name := ctx.Param("name")
		room, err := models.NewRoomModel().GetRoomByName(name)
		if err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusNotFound)
			return
		}

		// Query for join and quit room
		if action, ok := ctx.GetQuery("action"); ok {
			if action == "join" {
				username := utils.GetSessionValue(user)["username"].(string)

				for _, member := range room.Members {
					if member.Username == username {
						flash.Set("danger", "You has already joined.")
						ctx.Redirect(http.StatusFound, "/room/"+room.Name)
						return
					}
				}

				models.NewRoomModel().AddMember(room.ID, username)

				flash.Set("success", "Success join to this room.")
				ctx.Redirect(http.StatusFound, "/room/"+room.Name)
				return
			}

			if action == "quit" {
				username := utils.GetSessionValue(user)["username"].(string)

				models.NewRoomModel().RemoveMember(room.ID, username)

				flash.Set("success", "Success quit.")
				ctx.Redirect(http.StatusFound, "/room/"+room.Name)
				return
			}
		}

		var mentor models.User
		mentor, _ = models.NewUserModel().GetUserByID(room.MentorID)

		defer flash.Delete()
		ctx.HTML(http.StatusOK, "room_detail", gin.H{
			"user":   user,
			"mentor": mentor,
			"room":   room,
			"flash":  flash,
		})
	}
}
