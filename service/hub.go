package service

import "sh.ieki.xyz/model"

func NewHub() *model.Hub {
	hub := model.Hub{}
	hub.Construct()
	return &hub
}
