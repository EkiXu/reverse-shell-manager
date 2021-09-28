package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
	Detail string      `json:"detail"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(statuCode int, code int, data interface{}, detial string, c *gin.Context) {
	// 开始时间
	c.JSON(statuCode, Response{
		code,
		data,
		detial,
	})
}

func Ok(c *gin.Context) {
	Result(http.StatusOK, SUCCESS, map[string]interface{}{}, "success", c)
}

func OkWithData(statuCode int, data interface{}, c *gin.Context) {
	Result(statuCode, SUCCESS, data, "success", c)
}

func OkWithDetail(detail string, c *gin.Context) {
	Result(http.StatusOK, SUCCESS, map[string]interface{}{}, detail, c)
}

func Fail(statuCode int, c *gin.Context) {
	Result(statuCode, ERROR, map[string]interface{}{}, "something wrong", c)
}

func FailWithDetail(statuscode int, detail string, c *gin.Context) {
	Result(statuscode, ERROR, map[string]interface{}{}, detail, c)
}

func FailWithCodeAndDetail(statuscode int, code int, detail string, c *gin.Context) {
	Result(http.StatusInternalServerError, code, map[string]interface{}{}, detail, c)
}
