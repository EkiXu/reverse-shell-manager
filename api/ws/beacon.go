package ws

import (
	"errors"

	"github.com/gin-gonic/gin"
	"sh.ieki.xyz/global"
	"sh.ieki.xyz/global/typing"
	"sh.ieki.xyz/middleware"
	"sh.ieki.xyz/model"
	"sh.ieki.xyz/util"
)

/*
func DemoWSAPI(c *gin.Context) {
	wsConn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.SERVER_LOG.Errorf("failed to set websocket upgrade: %+v", err)
		return
	}
	defer wsConn.Close()

	listenerName, exist := c.GetQuery("name")

	if exist == false {
		return
	}

	listener, err := util.GetListener(global.SERVER_LISTENER_LIST, model.Listener{Name: listenerName})

	if err != nil {
		global.SERVER_LOG.Error(err)
		return
	}

	//listener.ServeWsConn(wsConn)

	//defer listener.StopServeWsConn(wsConn)

	for {
		var wsMsg typing.ShData
		err := wsConn.ReadJSON(&wsMsg)
		if err != nil {
			global.SERVER_LOG.Debug(err)
			break
		}
		global.SERVER_LOG.Debugf("ReceiveInput %v", wsMsg)

		if wsMsg.Type == "cmd" {
			err := listener.ServeWsCmdInput(wsConn, wsMsg.Data)
			if err != nil {
				global.SERVER_LOG.Debug(err)
				wsConn.WriteJSON(typing.ShData{Type: "error", Data: err.Error()})
			}
		}
	}
}*/

func BeaconWSAPI(c *gin.Context) {
	wsConn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		global.SERVER_LOG.Errorf("failed to set websocket upgrade: %+v", err)
		return
	}
	defer func() {
		if err != nil {
			wsConn.WriteJSON(typing.ShData{Type: "error", Data: err.Error()})
			wsConn.Close()
		}
	}()

	beaconName, isset := c.GetQuery("name")

	if !isset {
		err = errors.New("username required")
		return
	}

	var wsMsg typing.ShData
	err = wsConn.ReadJSON(&wsMsg)

	if wsMsg.Type != "auth" {
		err = errors.New("authentication required")
		return
	}

	jwt := middleware.NewJWT()
	_, err = jwt.ParseToken(wsMsg.Data)

	if err != nil {
		return
	}

	beacon, err := util.GetBeacon(global.SERVER_BEACON_LIST, model.Beacon{Name: beaconName})

	if err != nil {
		return
	}

	for {

		err := wsConn.ReadJSON(&wsMsg)
		if err != nil {
			global.SERVER_LOG.Debug(err)
			break
		}
		//global.SERVER_LOG.Debugf("ReceiveInput %v", wsMsg)

		if wsMsg.Type == "cmd" {
			err := beacon.ServeShDataInput(wsMsg, wsConn)
			if err != nil {
				global.SERVER_LOG.Debug(err)
				wsConn.WriteJSON(typing.ShData{Type: "error", Data: err.Error()})
			}
		}
	}
}
