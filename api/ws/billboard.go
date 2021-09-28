package ws

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"sh.ieki.xyz/global"
	"sh.ieki.xyz/middleware"
	"sh.ieki.xyz/model"

	"sh.ieki.xyz/global/typing"
)

func BillboardWSAPI(c *gin.Context) {
	wsConn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.SERVER_LOG.Errorf("failed to set websocket upgrade: %+v", err)
		return
	}

	defer func() {
		if err != nil {
			global.SERVER_LOG.Debugf("Error %+v", err)
			wsConn.WriteJSON(typing.WSData{Timestamp: model.GetNowTimeStamp(), Sender: "server", Type: "user", Detail: err.Error()})
			wsConn.Close()
		}
	}()

	username, isset := c.GetQuery("name")

	if !isset {
		err = errors.New("name required")
		return
	}

	if err != nil {
		global.SERVER_LOG.Error(err)
		return
	}

	var wsMsg typing.ShData
	err = wsConn.ReadJSON(&wsMsg)

	if err != nil {
		return
	}

	if wsMsg.Type != "auth" {
		err = errors.New("authentication required")
		return
	}

	jwt := middleware.NewJWT()
	claim, err := jwt.ParseToken(wsMsg.Data)

	if err != nil {
		return
	}

	if username != claim.Name {
		return
	}

	//listener.ServeWsConn(wsConn)

	//defer listener.StopServeWsConn(wsConn)

	user := &model.User{}

	user.Construct(global.SERVER_WS_HUB.(*model.Hub), username, wsConn)

	user.Hub.Register <- user

	user.Hub.BroadcastWSData(typing.WSData{Sender: "server", Type: "user", Detail: fmt.Sprintf("User %s joined in", username)})

	go user.ReadPump()
	go user.WritePump()
}
