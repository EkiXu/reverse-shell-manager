package initialize

import (
	"sh.ieki.xyz/global"
	"sh.ieki.xyz/service"
)

func HubInit() {
	global.SERVER_WS_HUB = service.NewHub()

	go global.SERVER_WS_HUB.Run()
}
