package controller

import (
	"pm-ssl-management/src/conf"
	"pm-ssl-management/src/log"
	"pm-ssl-management/src/util"

	"github.com/gin-gonic/gin"
)

var loginPrivate = ""

func firstLogin() string {
	private, public := util.GenerateRSAKeyPair()
	if private == "" || public == "" {
		return ""
	}
	loginPrivate = private
	return public
}

func Login(ctx *gin.Context) {
	var req map[string]string
	ctx.BindJSON(&req)

	if req["password"] == "" {
		public := firstLogin()

		if public == "" {
			ctx.JSON(200, util.SendErrModel(1))
			return
		}
		ctx.JSON(200, util.SendSusModel(public))

		return
	}

	cipherText := util.Encrypt(req["password"], loginPrivate)
	if cipherText == "" {
		ctx.JSON(200, util.SendErrModel(4))
		log.ComLoggerClient.Error("用户名或密码错误")
	}

	if cipherText != conf.GetConf("password") {
		ctx.JSON(200, util.SendErrModel(4))
		log.ComLoggerClient.Error("用户名或密码错误")
		return
	}

	ctx.JSON(200, util.SendSusModel(map[string]string{
		"token": util.GenerateToken(),
	}))
}
