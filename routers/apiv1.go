package routers

import (
	"kafekoding/handlers/apiv1"

	"github.com/gin-gonic/gin"
)

func APIV1Router(r *gin.Engine) {
	v1 := r.Group("/api/v1")

	// Room
	room := apiv1.NewAPIV1Room()
	v1.GET("/room", room.GetAllRooms())
	v1.GET("/room/:name", room.GetRoomByName())
}
