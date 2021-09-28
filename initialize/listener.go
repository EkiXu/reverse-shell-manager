package initialize

import (
	"container/list"

	"sh.ieki.xyz/global"
)

func ListenerInit() {
	global.SERVER_LISTENER_LIST = list.New()
}
