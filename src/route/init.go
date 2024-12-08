package route

import (
	"pm-ssl-management/src/conf"
	"pm-ssl-management/src/log"
	"pm-ssl-management/src/middleware"
	"pm-ssl-management/src/util"

	"github.com/gin-gonic/gin"
)

// {功能} 设置中间件
// {参数} 路由实例
// {返回} 无
func setMiddleware(r *gin.Engine) {
	r.Use(middleware.Recover)
	if util.IsDEV() {
		r.Use(middleware.UseCors)
	}
	r.Use(middleware.UseLog)
}

// {功能} 生成路由
// {参数} 无
// {返回} 无
func Init() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.ForwardedByClientIP = true

	setMiddleware(router)

	generateCertificate(router)
	generateLogin(router)

	port := conf.GetConf("port", "4006")

	log.ComLoggerFmt("服务启动端口: ", port)

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
