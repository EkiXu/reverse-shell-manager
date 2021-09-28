package router

import (
	"github.com/gin-gonic/gin"
	ws "sh.ieki.xyz/api/ws"
)

func InitBillboardRouter(Router *gin.RouterGroup) {
	ListenerRouter := Router.Group("billboard")
	{
		ListenerRouter.GET("/ws/", ws.BillboardWSAPI)
	}
}
