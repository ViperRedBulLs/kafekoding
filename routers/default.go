package routers

import (
	"kafekoding/handlers"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {
	r.GET("/login", handlers.Login())
	r.POST("/login", handlers.Login())
	r.GET("/register", handlers.Register())
	r.POST("/register", handlers.Register())
	r.GET("/logout", handlers.Logout())
}

func RoomRouter(r *gin.Engine) {
	r.GET("/", handlers.RoomIndex())
	r.POST("/", handlers.RoomIndex())
	r.GET("/room/:name", handlers.RoomDetail())
}

func AdminRouter(r *gin.RouterGroup) {
	r.GET("/login", handlers.AdminLogin())
	r.POST("/login", handlers.AdminLogin())
	r.GET("/detail", handlers.AdminDetail)
	r.GET("/", handlers.AdminIndex())
}

func ChatRouter(r *gin.Engine) {
	r.GET("/chat", handlers.ChatIndex)
	r.GET("/chat/:roomID", handlers.ChatDetail)
	r.POST("/chat/:roomID", handlers.ChatDetail)
}
