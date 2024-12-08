package middleware

import (
	"github.com/gin-gonic/gin"
)

// {功能} 跨域处理中间件
// {参数} 无
// {返回} 无
func UseCors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
	}
}
