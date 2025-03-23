package api

import (
	"gin/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func (con UserController) Index(c *gin.Context) {
	a := &models.Article{
		Title:   "首页页面",
		Content: "首页详情",
		Score:   91,
	}
	c.HTML(http.StatusOK, "api/index.html", gin.H{
		"title":     "前台首页",
		"news":      a,
		"hobby":     []string{"吃饭", "睡觉", "写代码"},
		"date":      time.Now().Unix(),
		"printInfo": "Building...",
	})
}
func (con UserController) News(c *gin.Context) {
	a := &models.Article{
		Title:   "api新闻页面",
		Content: "api新闻详情",
	}
	c.HTML(http.StatusOK, "api/news.html", gin.H{
		"title": "api新闻",
		"news":  a,
	})

}
