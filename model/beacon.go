package model

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sh.ieki.xyz/config"
	"sh.ieki.xyz/global"
	"sh.ieki.xyz/global/typing"
)

type Beacon struct {
	UUID        uuid.UUID    `json:"uuid"`
	Name        string       `json:"name"`
	shConn      *net.TCPConn `json:"-"`
	Listener    *Listener    `json:"-"`
	send        chan []byte  `json:"-"`
	receive     chan []byte  `json:"-"`
	stoppedChan chan bool    `json:"-"` //becon 掉线
	lck         sync.Mutex   `json:"-"` //一次只能执行一个用户的命令
}

func (b *Beacon) Construct(shConn *net.TCPConn) {
	bid := uuid.New()
	b.UUID = bid
	b.Name = bid.String()
	b.shConn = shConn
	b.send = make(chan []byte, 1024)
	b.receive = make(chan []byte, 1024)
	b.stoppedChan = make(chan bool, 1)
}

//beacon写
func (b *Beacon) WritePump() {
	//ticker := time.NewTicker(config.BC_Ping_Period)
	defer func() {
		//ticker.Stop()
		b.shConn.Close()
	}()
	for {
		select {
		case message, ok := <-b.send:
			b.shConn.SetWriteDeadline(time.Now().Add(config.BC_Write_Wait))
			if !ok {
				// The Beacon closed the channel.
				return
			}
			//global.SERVER_LOG.Debugf("Send message %v", string(message))
			b.shConn.Write(message)
			/*
				case <-ticker.C:
					b.shConn.SetWriteDeadline(time.Now().Add(config.BC_Write_Wait))
					if err := b.shConn.WriteMessage(websocket.PingMessage, nil); err != nil {
						return
					}
				}*/
		default:
			continue
		}
	}
}

func (b *Beacon) ReadPump() {

	defer func() {
		global.SERVER_WS_HUB.BroadcastWSData(typing.WSData{Timestamp: GetNowTimeStamp(), Sender: "server", Type: "beacon", Detail: fmt.Sprintf("Beacon %s go offline", b.Name)})
		b.shConn.Close()
		DeleteBeacon(global.SERVER_BEACON_LIST, b)
	}()

	//b.shConn.SetReadDeadline(time.Now().Add(config.WS_Pong_Wait))
	//b.shConn.SetPongHandler(func(string) error { u.wsConn.SetReadDeadline(time.Now().Add(config.WS_Pong_Wait)); return nil })
	buf := make([]byte, 2048)
	for {

		n, err := b.shConn.Read(buf)
		//global.SERVER_LOG.Debugf("Received message %v", buf[0:n])
		if err != nil {
			global.SERVER_LOG.Errorf("Becon %s error: %v", b.Name, err)
			break
		}
		if n > 0 {
			//global.SERVER_LOG.Debugf("Send message %v to WsConn", string(buf[0:n]))
			b.receive <- buf[0:n]
		}
	}
}

func (b *Beacon) ServeShDataInput(shData typing.ShData, wsConn *websocket.Conn) error {
	b.lck.Lock()
	defer b.lck.Unlock()

	//global.SERVER_LOG.Debugf("ReceiveshDataInput %+v", shData)

	switch shData.Type {
	case "cmd":
		b.send <- []byte(shData.Data)
	}

	//global.SERVER_LOG.Debugf("Sended")

	timeout := time.NewTimer(config.BC_CMD_Wait)

	receivedOnce := false
	for {
		select {
		case message, ok := <-b.receive:
			//write combine output bytes into websocket response
			//global.SERVER_LOG.Debugf("WsConn receoved %v", string(message))
			if !ok {
				return errors.New("Beacon closed")
			}
			wsConn.WriteJSON(typing.ShData{Type: "cmd", Data: string(message)})
			receivedOnce = true
			timeout.Stop()
			timeout.Reset(config.BC_CMD_Reset_Wait)

		case <-timeout.C:
			if !receivedOnce {
				return errors.New("timeout")
			}
			return nil
		case <-b.stoppedChan:
			return errors.New("beacon closed")
		}
	}
}

func DeleteBeacon(beaconList *list.List, targetBeacon *Beacon) error {
	var beacon *Beacon
	// global.SERVER_LOG.Debugf("finding beacon,target beacon: %+v", targetBeacon)
	for element := beaconList.Front(); element != nil; element = element.Next() {
		beacon = element.Value.(*Beacon)

		// global.SERVER_LOG.Debugf("finding beacon,now beacon: %+v", *beacon)
		if (*beacon).Name == targetBeacon.Name || (*beacon).UUID == targetBeacon.UUID {
			beaconList.Remove(element)
			return nil
		}
	}
	return errors.New("beacon not found")
}

func GetNowTimeStamp() int64 {
	return time.Now().UnixNano() / 1e6
}
