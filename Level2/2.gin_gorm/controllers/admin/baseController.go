package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

func (con BaseController) success(c *gin.Context) {
	value, isExist := c.Get("username")
	cookie, err := c.Cookie("cookieName")

	if isExist {
		if err == nil {
			c.String(http.StatusOK, "中间件传value=%v, 浏览器cookie=%v", value, cookie)
		} else {
			c.String(http.StatusOK, "中间件与控制器共享数据: value=%v", value)
		}

	} else {
		if err == nil {
			c.String(http.StatusOK, "cookie=%v", cookie)
		} else {
			c.String(http.StatusOK, "请求成功")
		}

	}

}

// func (con BaseController) error(c *gin.Context) {
// 	c.String(http.StatusOK, "失败")
// }
