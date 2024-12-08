package middleware

import (
	"pm-ssl-management/src/util"

	"github.com/gin-gonic/gin"
)

// {功能} 跨域处理中间件
// {参数} 无
// {返回} 无
func UseToken(c *gin.Context) {
	// 获取token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(200, util.SendErrModel(5))
		c.Abort()
		return
	}
	// 验证token
	if !util.ValidateToken(token) {
		c.JSON(200, util.SendErrModel(5))
		c.Abort()
		return
	}
	c.Next()
}
