package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"sh.ieki.xyz/global"
	"sh.ieki.xyz/global/response"
	"sh.ieki.xyz/global/typing"
	"sh.ieki.xyz/service"
	"sh.ieki.xyz/util"
)

// @Tags Listener
// @Summary 添加
// @Produce  application/json
// @Param data body request.addListenerStruct true "添加监听器接口"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Register Successfully"}"
// @Router /listener/ [post]
func AddListenerAPI(c *gin.Context) {
	var request typing.AddListenerRequestStruct
	err := c.ShouldBindJSON(&request)

	if err != nil {
		response.FailWithDetail(http.StatusBadRequest, err.Error(), c)
		return
	}

	newListener, err := service.AddListener(request.Name, request.LHOST, request.LPORT)
	if err != nil {
		response.FailWithDetail(http.StatusBadRequest, err.Error(), c)
		global.SERVER_LOG.Error(err)
		return
	}
	global.SERVER_WS_HUB.BroadcastWSData(typing.WSData{Timestamp: time.Now().UnixNano() / 1e6, Sender: "server", Type: "listener", Data: *newListener, Detail: fmt.Sprintf("Start Listening At %s:%d", newListener.Host, newListener.Port)})
	response.OkWithData(http.StatusCreated, *newListener, c)
}

func ListenerListAPI(c *gin.Context) {
	listenerList, err := util.GetListenerList(global.SERVER_LISTENER_LIST)
	if err != nil {
		response.FailWithDetail(http.StatusInternalServerError, err.Error(), c)
		return
	}
	response.OkWithData(http.StatusOK, listenerList, c)
	return
}
