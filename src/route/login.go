package route

import (
	"pm-ssl-management/src/controller"

	"github.com/gin-gonic/gin"
)

func generateLogin(router *gin.Engine) {
	certificate := router.Group("login")

	certificate.POST("", controller.Login)
}
