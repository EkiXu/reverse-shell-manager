package model

import (
	"fmt"
	"net"

	"sh.ieki.xyz/global"
	"sh.ieki.xyz/global/typing"
)

//反弹shell监听器
type Listener struct {
	Name        string       `json:"name"`
	Host        string       `json:"lhost"`
	Port        int          `json:"lport"`
	Closed      bool         `json:"closed"`
	socket      net.Listener //listener 启动的socket
	stoppedChan chan bool    //停止TCP socket监听
}

// Listen for incoming connections
func (listener *Listener) serve() {
	for {
		conn, err := listener.socket.Accept()
		// FIXME I'd like to detect "use of closed network connection" here
		// FIXME but it isn't exported from net
		if err != nil {
			global.SERVER_LOG.Errorf("Accept failed: %v", err)
			break
		}
		//Initial Read
		buf := make([]byte, 4096)
		n, _ := conn.(*net.TCPConn).Read(buf)
		if n > 0 {
			global.SERVER_LOG.Debugf("Shell received: %s", string(buf[0:n]))
		}

		beacon := &Beacon{}

		beacon.Construct(conn.(*net.TCPConn))

		_ = global.SERVER_WS_HUB.BroadcastWSData(typing.WSData{Sender: "server", Type: "beacon", Data: beacon, Detail: fmt.Sprintf("Shell received: %s", string(buf[0:n]))})

		global.SERVER_LOG.Debugf("before pushback beacon list:%+v and push %v", global.SERVER_BEACON_LIST, beacon)
		global.SERVER_BEACON_LIST.PushBack(beacon)
		global.SERVER_LOG.Debugf("after pushback beacon list:%v", global.SERVER_BEACON_LIST)
		go beacon.WritePump()
		go beacon.ReadPump()
	}
	listener.stoppedChan <- true
}

// Stop the server by closing the listening listen
func (listener *Listener) Stop() {
	listener.socket.Close()
	listener.Closed = true
	<-listener.stoppedChan
}

func (listener *Listener) Start() error {
	addr := fmt.Sprintf("%s:%d", listener.Host, listener.Port)
	var err error
	listener.socket, err = net.Listen("tcp", addr)
	listener.stoppedChan = make(chan bool, 1)
	if err != nil {
		global.SERVER_LOG.Error("tcp listener error")
		return err
	}
	global.SERVER_LOG.Debugf("listener %s start listen at %s", listener.Name, addr)
	listener.Closed = false
	go listener.serve()
	return nil
}

/*
func (listener *Listener) send(remote *net.TCPConn) {
	defer remote.Close()
	for {
		select {
		case <-listener.receivedChan:
			//write combine output bytes into websocket response
			//write combine output bytes into websocket response
			//global.SERVER_LOG.Debugf("Enter send with %+v", listener, listener.shInputBuffer)
			buf := make([]byte, 4096)
			n, err := listener.shInputBuffer.Read(buf)
			if err == nil && n > 0 {
				global.SERVER_LOG.Debugf("Send %+v to shConn", string(buf[0:n]))
				if _, err := remote.Write(buf[0:n]); err != nil {
					global.SERVER_LOG.Error("data send to ShConn failed")
				}
			}
		case <-listener.stoppedChan:
			return
		}
	}
}

func (listener *Listener) respond(remote *net.TCPConn) {
	defer remote.Close()
	tick := time.NewTicker(time.Millisecond * time.Duration(120))
	buf := make([]byte, 4096)
	for {
		select {
		case <-tick.C:
			//write combine output bytes into websocket response
			n, _ := remote.Read(buf)
			if n > 0 {
				global.SERVER_LOG.Debugf("ReceiveWsMsg len:%d content:%s", n, string(buf[0:n]))
				listener.shOutputBuffer.Write(buf[0:n])
			}
		case <-listener.stoppedChan:
			return
		}
	}
}



func (listener *Listener) SendWsCmdInput(cmd string) (string, error) {
	global.SERVER_LOG.Debug("enter here")
	if listener.Closed {
		return "", errors.New("listener closed")
	}

	//一次处理一个命令0
	listener.lck.Lock()
	defer listener.lck.Unlock()

	_, err := listener.shInputBuffer.Write([]byte(cmd))

	if err != nil {
		listener.shInputBuffer.Reset()
		listener.shOutputBuffer.Reset()
		return "", err
	}

	listener.receivedChan <- true

	select {
	case <-listener.doneChan:
		buf := make([]byte, 4096)
		n, err := listener.shOutputBuffer.Read(buf)
		if err != nil {
			return "", err
		}
		return string(buf[0:n]), nil
	case <-time.After(time.Second * 5):
		return "", errors.New("cmd timeout")
	case <-listener.stoppedChan:
		return "", errors.New("listener stopped")
	}
}

func (listener *Listener) ServeWsCmdInput(wsConn *websocket.Conn, cmd string) error {
	global.SERVER_LOG.Debug("enter here")
	if listener.Closed {
		return errors.New("listener closed")
	}

	//一次处理一个命令0
	listener.lck.Lock()
	defer listener.lck.Unlock()

	_, err := listener.shInputBuffer.Write([]byte(cmd))

	if err != nil {
		return err
	}

	listener.receivedChan <- true

	tick := time.NewTicker(time.Millisecond * time.Duration(121))
	timeout := time.NewTimer(time.Duration(time.Second * 5))
	buf := make([]byte, 4096)
	receivedOnce := false
	for {
		select {
		case <-tick.C:
			//write combine output bytes into websocket response
			n, _ := listener.shOutputBuffer.Read(buf)
			if n > 0 {
				wsConn.WriteJSON(typing.ShData{Type: "cmd", Data: string(buf[0:n])})
				receivedOnce = true
				timeout.Stop()
				timeout.Reset(time.Millisecond * time.Duration(361))
			}
		case <-timeout.C:
			if !receivedOnce {
				return errors.New("timeout")
			}
			return nil
		case <-listener.stoppedChan:
			return errors.New("listener closed")
		}
	}
}
*/
