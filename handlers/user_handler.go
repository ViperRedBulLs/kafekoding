package handlers

import (
	"kafekoding/models"
	"kafekoding/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Login handler for user login
func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		if user := session.Get("user"); user != nil {
			ctx.Redirect(http.StatusFound, "/")
			return
		}

		// Method GET
		if ctx.Request.Method == http.MethodGet {
			defer flash.Delete()

			ctx.HTML(http.StatusOK, "login", gin.H{
				"flash": flash,
			})
			return
		}

		// Method POST
		if ctx.Request.Method == http.MethodPost {
			userName := ctx.PostForm("username")
			password := ctx.PostForm("password")

			user, err := models.NewUserModel().GetUserByUsername(userName)
			if err != nil {
				log.Println(err.Error())
				flash.Set("danger", "Password is wrong.")
				ctx.Redirect(http.StatusFound, "/login")
				return
			}

			if !utils.DecryptionPassword(user.Password, password) {
				flash.Set("danger", "Password is wrong.")
				ctx.Redirect(http.StatusFound, "/login")
				return
			}

			user.IsOnline = true
			user.LastLogin = time.Now()
			err = models.GetDB.Save(&user).Error
			if err != nil {
				http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
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

			ctx.Redirect(http.StatusFound, "/")
			return
		}
	}
}

// Register handler for registering a user
func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		if user := session.Get("user"); user != nil {
			ctx.Redirect(http.StatusFound, "/")
			return
		}

		// Method GET
		if ctx.Request.Method == http.MethodGet {
			defer flash.Delete()

			ctx.HTML(http.StatusOK, "register", gin.H{
				"flash": flash,
			})
			return
		}

		// Method POST
		if ctx.Request.Method == http.MethodPost {
			firstName := ctx.PostForm("first_name")
			lastName := ctx.PostForm("last_name")
			username := ctx.PostForm("username")
			email := ctx.PostForm("email")
			password := ctx.PostForm("password")

			user := models.User{
				FirstName: firstName,
				LastName:  lastName,
				Username:  username,
				Email:     email,
				Password:  password,
			}

			err := models.NewUserModel().CreateUser(&user)
			if err != nil {
				flash.Set("danger", err.Error())
				ctx.Redirect(http.StatusFound, "/register")
				return
			}

			flash.Set("success", "Successfully to create user by username: "+username)
			ctx.Redirect(http.StatusFound, "/login")
			return
		}
	}
}

// Logout handler
func Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		userSession := session.Get("user")

		sessionValue := utils.GetSessionValue(userSession)
		user, err := models.NewUserModel().GetUserByID(sessionValue["id"].(uint))
		if err != nil {
			log.Println(err)
			return
		}
		user.IsOnline = false
		models.GetDB.Save(&user)

		session.Delete("user")
		session.Clear()
		if err := session.Save(); err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx.Redirect(http.StatusFound, "/login")
	}
}
