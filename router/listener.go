package router

import (
	"github.com/gin-gonic/gin"
	api "sh.ieki.xyz/api"
	"sh.ieki.xyz/middleware"
)

func InitListenerRouter(Router *gin.RouterGroup) {
	ListenerRouter := Router.Group("listener").Use(middleware.JWTAuth())
	{
		ListenerRouter.POST("/", api.AddListenerAPI)
		ListenerRouter.GET("/", api.ListenerListAPI)
	}
}
