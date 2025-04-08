package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")//允许的来源（这里设置为 *，表示允许所有来源）
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")//允许的 HTTP 方法
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, authCode, token, Content-Type, Accept, Authorization")//允许的请求头字段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")//是否允许携带凭据
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
