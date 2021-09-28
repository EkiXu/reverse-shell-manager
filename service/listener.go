package service

import (
	"errors"

	"sh.ieki.xyz/global"
	"sh.ieki.xyz/model"
	"sh.ieki.xyz/util"
)

func AddListener(name string, host string, port int) (*model.Listener, error) {
	listener := model.Listener{
		Name:   name,
		Host:   host,
		Port:   port,
		Closed: true,
	}

	_, err := util.GetListener(global.SERVER_LISTENER_LIST, listener)

	if err == nil {
		return nil, errors.New("listener already created")
	}

	err = listener.Start()

	if err != nil {
		return nil, err
	}

	global.SERVER_LISTENER_LIST.PushBack(&listener)

	return &listener, nil
}
