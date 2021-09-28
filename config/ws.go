package config

import "time"

const (
	// Time allowed to write a message to the peer.
	WS_Write_Wait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	WS_Pong_Wait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	WS_Ping_Period = (WS_Pong_Wait * 9) / 10

	// Maximum message size allowed from peer.
	WS_Max_Message_Size = 512
)
