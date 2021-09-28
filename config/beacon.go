package config

import "time"

const (
	// Time allowed to write a message to the peer.
	BC_Write_Wait = 10 * time.Second

	BC_CMD_Wait = 10 * time.Second

	BC_CMD_Reset_Wait = 1 * time.Second

	// Time allowed to read the next pong message from the peer.
	BC_Pong_Wait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	BC_Ping_Period = (BC_Pong_Wait * 9) / 10

	// Maximum message size allowed from peer.
	BC_Max_Message_Size = 2048
)
