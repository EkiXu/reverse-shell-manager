package router

import (
	"github.com/gin-gonic/gin"
	"sh.ieki.xyz/api"
)

func InitTokenRouter(Router *gin.RouterGroup) {
	TokenRouter := Router.Group("token")
	{
		TokenRouter.POST("/", api.LoginAPI)
	}
}
