package core

import (
	"fmt"
	"net/http"
	"time"

	"sh.ieki.xyz/global"
	"sh.ieki.xyz/initialize"
)

func RunServer() {
	Router := initialize.RouterInit()
	initialize.ListenerInit()
	initialize.HubInit()
	initialize.BeaconListInit()

	address := fmt.Sprintf(":%d", global.SERVER_CONFIG.System.Addr)
	s := &http.Server{
		Addr:           address,
		Handler:        Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.SERVER_LOG.Debug("server run success on ", address)

	global.SERVER_LOG.Error(s.ListenAndServe())
}
