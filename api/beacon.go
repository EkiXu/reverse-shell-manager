package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sh.ieki.xyz/global"
	"sh.ieki.xyz/global/response"
	"sh.ieki.xyz/util"
)

func BeaconListAPI(c *gin.Context) {
	beaconList, err := util.GetBeaconList(global.SERVER_BEACON_LIST)
	if err != nil {
		response.FailWithDetail(http.StatusInternalServerError, err.Error(), c)
		return
	}
	response.OkWithData(http.StatusOK, beaconList, c)
	return
}
