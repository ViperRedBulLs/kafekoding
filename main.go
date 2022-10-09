package main

import (
	"encoding/gob"
	"kafekoding/handlers/backends"
	"kafekoding/models"
	"kafekoding/routers"
	"log"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	gob.Register(time.Time{})
	gob.Register(map[string]interface{}{})
}

func main() {
	r := gin.Default()
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Println(err.Error())
		return
	}
	r.Static("/static", "./static")
	r.Static("/media", "./media")
	r.HTMLRender = createMyRender()
	r.Use(sessions.Sessions("ginsessionID", models.GetStore()))

	routers.RoomRouter(r)
	routers.UserRouter(r)

	// API
	routers.APIV1Router(r)

	// chat
	routers.ChatRouter(r)

	// Admin
	admin := r.Group("/admin")
	routers.AdminRouter(admin)

	hub := backends.NewHub()
	go hub.Run()
	// Websocket Chat
	r.GET("/backends/ws/chat", func(ctx *gin.Context) {
		room := ctx.Query("room")
		backends.ServeWS(hub, ctx.Writer, ctx.Request, room)
	})

	err = r.Run(":8000")
	if err != nil {
		log.Println(err)
		return
	}
}
