package initialize

import (
	"container/list"

	"sh.ieki.xyz/global"
)

func BeaconListInit() {
	global.SERVER_BEACON_LIST = list.New()
}
