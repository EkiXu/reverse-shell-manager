package api

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"sh.ieki.xyz/global"
	"sh.ieki.xyz/global/response"
	"sh.ieki.xyz/global/typing"
	"sh.ieki.xyz/middleware"
)

// @Tags Listener
// @Summary 添加
// @Produce  application/json
// @Param data body request.addListenerStruct true "添加监听器接口"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Register Successfully"}"
// @Router /listener/ [post]
func LoginAPI(c *gin.Context) {
	var request typing.LoginRequestStruct
	err := c.ShouldBindJSON(&request)

	if err != nil {
		response.FailWithDetail(http.StatusBadRequest, err.Error(), c)
		return
	}

	has := md5.Sum([]byte(request.Password))
	hasHex := fmt.Sprintf("%x", has)

	if hasHex == global.SERVER_CONFIG.Auth.PasswordHash {
		tokenNext(c, request.Username)
		return
	}

	response.FailWithDetail(http.StatusForbidden, "Wrong Password", c)
}

// 登录以后签发jwt
func tokenNext(c *gin.Context, username string) {
	j := &middleware.JWT{
		SigningKey: []byte(global.SERVER_CONFIG.Auth.JWTKey), // 唯一签名
	}
	claims := typing.CustomClaims{
		Name: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,       // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, // 过期时间 一周
			Issuer:    "Xuanji",                       // 签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		response.FailWithDetail(http.StatusInternalServerError, "Failed to get token", c)
		return
	}
	response.OkWithData(http.StatusCreated, typing.LoginResponseStruct{
		Username:    username,
		AccessToken: token,
		ExpiresAt:   claims.StandardClaims.ExpiresAt * 1000,
	}, c)
	return
}
