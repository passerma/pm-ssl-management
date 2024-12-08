package middleware

import (
	"net/http"
	"pm-ssl-management/src/log"
	"pm-ssl-management/src/util"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			stackStr := string(debug.Stack())
			log.ComLoggerClient.Error(errorToString(r)+"\n", stackStr)
			c.JSON(http.StatusOK, util.SendErrModel(1))
			c.Abort()
		}
	}()
	c.Next()
}

// recover错误，转string
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
