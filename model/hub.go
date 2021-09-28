package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"sh.ieki.xyz/config"
	"sh.ieki.xyz/global"

	"sh.ieki.xyz/global/typing"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type User struct {
	Hub    *Hub            `json:"-"`
	Name   string          `json:"name"`
	wsConn *websocket.Conn `json:"-"`
	send   chan []byte     `json:"-"`
}

func (u *User) Construct(hub *Hub, name string, wsConn *websocket.Conn) {
	u.Hub = hub
	u.Name = name
	u.wsConn = wsConn
	u.send = make(chan []byte, 256)
}

//用户读管道 从客户端读
func (u *User) ReadPump() {
	defer func() {
		u.Hub.Unregister <- u
		u.Hub.BroadcastWSData(typing.WSData{Sender: "server", Type: "user", Detail: fmt.Sprintf("User %s go offline", u.Name)})
		u.wsConn.Close()
	}()
	u.wsConn.SetReadLimit(config.WS_Max_Message_Size)
	u.wsConn.SetReadDeadline(time.Now().Add(config.WS_Pong_Wait))
	u.wsConn.SetPongHandler(func(string) error { u.wsConn.SetReadDeadline(time.Now().Add(config.WS_Pong_Wait)); return nil })
	for {
		_, rawMessage, err := u.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				global.SERVER_LOG.Errorf("error: %v", err)
			}
			break
		}
		var wsData typing.WSData
		err = json.Unmarshal(rawMessage, &wsData)
		if err != nil {
			global.SERVER_LOG.Debugf("error Unmarshal %+v", rawMessage)
			continue
		}
		wsData.Sender = u.Name

		u.Hub.BroadcastWSData(wsData)
	}
}

//用户写通道 向客户端写
func (u *User) WritePump() {
	ticker := time.NewTicker(config.WS_Ping_Period)
	defer func() {
		ticker.Stop()
		u.wsConn.Close()
	}()
	for {
		select {
		case message, ok := <-u.send:
			u.wsConn.SetWriteDeadline(time.Now().Add(config.WS_Write_Wait))
			if !ok {
				// The hub closed the channel.
				u.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := u.wsConn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(u.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-u.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			u.wsConn.SetWriteDeadline(time.Now().Add(config.WS_Write_Wait))
			if err := u.wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

//公共频道
type Hub struct {
	// Registered clients.
	Clients map[*User]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	Register chan *User

	// Unregister requests from clients.
	Unregister chan *User
}

func (h *Hub) Construct() {
	h.broadcast = make(chan []byte)
	h.Register = make(chan *User)
	h.Unregister = make(chan *User)
	h.Clients = make(map[*User]bool)
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastWSData(wsData typing.WSData) error {
	wsData.Timestamp = GetNowTimeStamp()
	global.SERVER_LOG.Debugf("broadcast %+v", wsData)
	rawData, err := json.Marshal(wsData)
	if err != nil {
		return err
	}
	h.broadcast <- rawData
	return nil
}
