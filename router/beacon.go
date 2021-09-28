package router

import (
	"github.com/gin-gonic/gin"
	"sh.ieki.xyz/api"
	ws "sh.ieki.xyz/api/ws"
	"sh.ieki.xyz/middleware"
)

func InitBeaconRouter(Router *gin.RouterGroup) {
	BeaconRouter := Router.Group("beacon")
	{
		BeaconRouter.GET("/ws/", ws.BeaconWSAPI)
		BeaconRouter.Use(middleware.JWTAuth())
		BeaconRouter.GET("/", api.BeaconListAPI)

	}
}
