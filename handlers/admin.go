package handlers

import (
	"kafekoding/models"
	"kafekoding/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AdminLogin for login admin
func AdminLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user := session.Get("user")

		if user != nil {
			if !utils.GetSessionValue(user)["isSuperuser"].(bool) || !utils.GetSessionValue(user)["isStaff"].(bool) {
				flash.Set("danger", "You not SuperUser, please login with superuser account.")
				ctx.HTML(http.StatusOK, "admin_login", gin.H{
					"flash": flash,
				})
				return
			}
		}

		if ctx.Request.Method == http.MethodGet {
			defer flash.Delete()

			ctx.HTML(http.StatusOK, "admin_login", gin.H{
				"flash": flash,
			})
			return
		}

		if ctx.Request.Method == http.MethodPost {
			username := ctx.PostForm("username")
			password := ctx.PostForm("password")

			user, err := models.NewUserModel().GetUserByUsername(username)
			if err != nil {
				flash.Set("danger", "Username or password is wrong")
				ctx.Redirect(http.StatusFound, "/admin/login")
				return
			}

			if !user.IsSuperuser || !user.IsStaff {
				flash.Set("danger", "Your account is not SuperUser")
				ctx.Redirect(http.StatusFound, "/admin/login")
				return
			}

			if !utils.DecryptionPassword(user.Password, password) {
				flash.Set("danger", "Username or password is wrong")
				ctx.Redirect(http.StatusFound, "/admin/login")
				return
			}

			userSession := map[string]interface{}{
				"id":          user.ID,
				"firstName":   user.FirstName,
				"lastName":    user.LastName,
				"fullName":    user.FirstName + " " + user.LastName,
				"username":    user.Username,
				"email":       user.Email,
				"isSuperuser": user.IsSuperuser,
				"isStaff":     user.IsStaff,
				"isActive":    user.IsActive,
				"isOnline":    user.IsOnline,
				"lastLogin":   user.LastLogin,
				"dateJoined":  user.DateJoined,
			}

			session.Set("user", userSession)
			if err = session.Save(); err != nil {
				http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx.Redirect(http.StatusFound, "/admin")
			return
		}
	}
}

// AdminIndex
func AdminIndex() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user := session.Get("user")
		if user == nil {
			flash.Set("danger", "Please login first before you accessing this page.")
			ctx.Redirect(http.StatusFound, "/admin/login")
			return
		}

		if !utils.GetSessionValue(user)["isSuperuser"].(bool) || !utils.GetSessionValue(user)["isStaff"].(bool) {
			flash.Set("danger", "Please login with SuperUser and Staff user.")
			ctx.Redirect(http.StatusFound, "/admin/login")
			return
		}

		if ctx.Request.Method == http.MethodGet {

			users := models.NewUserModel().GetAllUsers()
			rooms := models.NewRoomModel().GetAllRooms()
			courses := models.NewCourseModel().GetAllCourses()
			comments := models.NewCommentModel().GetAllComments()

			var countUser int
			var countRoom int
			var countCourse int
			var countComment int
			for range comments {
				countComment += 1
			}
			for range users {
				countUser += 1
			}
			for range courses {
				countCourse += 1
			}
			for range rooms {
				countRoom += 1
			}

			ctx.HTML(http.StatusOK, "admin_index", gin.H{
				"user":         user,
				"rooms":        rooms,
				"users":        users,
				"courses":      courses,
				"countCourse":  countCourse,
				"countUser":    countUser,
				"countRoom":    countRoom,
				"countComment": countComment,
			})
		}
	}
}

func AdminDetail(ctx *gin.Context) {
	session := sessions.Default(ctx)

	user := session.Get("user")
	if user == nil {
		flash.Set("danger", "Please login first before you accessing this page.")
		ctx.Redirect(http.StatusFound, "/admin/login")
		return
	}

	if !utils.GetSessionValue(user)["isSuperuser"].(bool) || !utils.GetSessionValue(user)["isStaff"].(bool) {
		flash.Set("danger", "Please login with SuperUser and Staff user.")
		ctx.Redirect(http.StatusFound, "/admin/login")
		return
	}

	if ctx.Request.Method == http.MethodGet {
		defer flash.Delete()

		var adminAction AdminAction
		if err := ctx.Bind(&adminAction); err == nil {
			if adminAction.TableName == "users" {
				// For user detail by passing query User's ID
				u, err := models.NewUserModel().GetUserByID(adminAction.FieldID)
				if err != nil {
					http.Error(ctx.Writer, "Not found", http.StatusNotFound)
					return
				}
				ctx.HTML(http.StatusOK, "admin_detail", gin.H{
					"user": user,
					"u":    u,
				})
				return
			}

			if adminAction.TableName == "rooms" {
				room, err := models.NewRoomModel().GetRoomByID(adminAction.FieldID)
				if err != nil {
					http.Error(ctx.Writer, "Not found", http.StatusNotFound)
					return
				}

				users := models.NewUserModel().GetAllUsers()

				ctx.HTML(http.StatusOK, "admin_detail", gin.H{
					"user":  user,
					"room":  room,
					"users": users,
				})
				return
			}
		} else if err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
