package route

import (
	"pm-ssl-management/src/controller"
	"pm-ssl-management/src/middleware"

	"github.com/gin-gonic/gin"
)

func generateCertificate(router *gin.Engine) {
	certificate := router.Group("certificate")

	certificate.Use(middleware.UseToken)

	certificate.POST("", controller.PostCertificate)
	certificate.GET("", controller.GetCertificate)
	certificate.DELETE("/:id", controller.DeleteCertificate)
	certificate.PUT("/:id", controller.PutCertificate)

	certificate.POST("/apply/:id", controller.PostCertificateApply)
	certificate.GET("/state/:id", controller.GetCertificateState)
}
