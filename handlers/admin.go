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
				defer flash.Delete()
				ctx.HTML(http.StatusOK, "admin_detail", gin.H{
					"user":  user,
					"u":     u,
					"flash": flash,
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

				defer flash.Delete()
				ctx.HTML(http.StatusOK, "admin_detail", gin.H{
					"user":  user,
					"room":  room,
					"users": users,
					"flash": flash,
				})
				return
			}
		} else if err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

// AdminPOSTUser handler for post user
func AdminPOSTUser(ctx *gin.Context) {
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

	id := ctx.PostForm("id")
	firstName := ctx.PostForm("first_name")
	lastName := ctx.PostForm("last_name")
	username := ctx.PostForm("username")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	isSuperuser := ctx.PostForm("is_superuser")
	isStaff := ctx.PostForm("is_staff")
	isActive := ctx.PostForm("is_active")
	// lastLogin := ctx.PostForm("last_login")
	// dateJoined := ctx.PostForm("date_joined")

	idInt, _ := strconv.Atoi(id)
	isSuperuserBool, _ := strconv.ParseBool(isSuperuser)
	isStaffBool, _ := strconv.ParseBool(isStaff)
	isActiveBool, _ := strconv.ParseBool(isActive)

	u, err := models.NewUserModel().GetUserByID(uint(idInt))
	if err != nil {
		return
	}

	u.FirstName = firstName
	u.LastName = lastName
	u.Username = username
	u.Email = email
	u.Password = password
	u.IsSuperuser = isSuperuserBool
	u.IsActive = isActiveBool
	u.IsStaff = isStaffBool

	err = models.GetDB.Save(&u).Error
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	flash.Set("success", "Successfully upgrade user by username: "+username)
	ctx.Redirect(http.StatusFound, "/admin/detail?table_name=users&field_id="+id)
}

// AdminPOSTRoom handler for post room
func AdminPOSTRoom(ctx *gin.Context) {
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

	id := ctx.PostForm("id")
	name := ctx.PostForm("name")
	mentorID := ctx.PostForm("mentor_id")
	members := ctx.PostForm("members")
	logo, err := ctx.FormFile("logo")
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	desc := ctx.PostForm("desc")
	content := ctx.PostForm("content")

	idInt, _ := strconv.Atoi(id)
	mentorIDInt, _ := strconv.Atoi(mentorID)
	membersInt, _ := strconv.Atoi(members)
	log.Println(membersInt)

	filename := logo.Filename

	u, _ := models.NewUserModel().GetUserByID(uint(membersInt))

	room, err := models.NewRoomModel().GetRoomByID(uint(idInt))
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusNotFound)
		return
	}

	room.Name = name
	room.MentorID = uint(mentorIDInt)
	room.Members = append(room.Members, &u)
	room.Logo = "/media/rooms/" + filename
	room.Desc = desc
	room.Content = content

	err = models.GetDB.Save(&room).Error
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = ctx.SaveUploadedFile(logo, "media/rooms/"+filename)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	flash.Set("success", "Successfully upgrade the room by ID: "+id)
	ctx.Redirect(http.StatusFound, "/admin/detail?table_name=rooms&field_id="+id)
}
