package global

import (
	"container/list"

	"sh.ieki.xyz/config"
	"sh.ieki.xyz/global/typing"

	oplogging "github.com/op/go-logging"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	SERVER_DB            *gorm.DB
	SERVER_CONFIG        config.Server
	SERVER_VP            *viper.Viper
	SERVER_LOG           *oplogging.Logger
	SERVER_LISTENER_LIST *list.List
	SERVER_BEACON_LIST   *list.List
	SERVER_WS_HUB        typing.HubInterface
)
