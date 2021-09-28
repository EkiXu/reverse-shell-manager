package initialize

import (
	"net/http"

	"sh.ieki.xyz/api"
	"sh.ieki.xyz/global"
	"sh.ieki.xyz/router"

	"github.com/gin-gonic/gin"
)

// 初始化总路由
func RouterInit() *gin.Engine {
	var r = gin.Default()

	r.LoadHTMLGlob("dist/*.html")
	r.Static("/static", "dist/static")

	//Revershell as service
	r.GET("/gen/:lhost/:lport", api.GenReverseShellPayloadAPI)

	// 方便统一添加路由组前缀 多服务器上线使用
	apiGroup := r.Group("/api/v1")
	router.InitListenerRouter(apiGroup)
	router.InitBillboardRouter(apiGroup)
	router.InitBeaconRouter(apiGroup)
	router.InitTokenRouter(apiGroup)
	global.SERVER_LOG.Info("router register success")

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	return r
}
