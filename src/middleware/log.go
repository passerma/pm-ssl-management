package middleware

import (
	"pm-ssl-management/src/log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// {功能} 访问日志中间件
// {参数} 无
// {返回} 无
func UseLog(c *gin.Context) {
	logger := log.AccessLoggerClient
	// 请求方式
	method := c.Request.Method
	// 请求路由
	url := c.Request.RequestURI
	// 请求IP
	ip := c.ClientIP()
	logger.WithFields(logrus.Fields{
		"ip":     ip,
		"method": method,
		"url":    url,
	}).Info(c.HandlerName())
}
