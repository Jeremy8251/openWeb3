package routers

import (
	"gin/controllers/api"

	"github.com/gin-gonic/gin"
)

func ApiRoutersInit(router *gin.Engine) {
	apiRouters := router.Group("/api")
	{
		apiRouters.GET("/news", api.UserController{}.News)
		// func(c *gin.Context) {
		// 	a := &model.Article{
		// 		Title:   "api新闻页面",
		// 		Content: "api新闻详情",
		// 	}
		// 	c.HTML(http.StatusOK, "api/news.html", gin.H{
		// 		"title": "api新闻",
		// 		"news":  a,
		// 	})
		// })

		apiRouters.GET("/", api.UserController{}.Index)
		// func(c *gin.Context) {
		// 	a := &model.Article{
		// 		Title:   "首页页面",
		// 		Content: "首页详情",
		// 		Score:   91,
		// 	}
		// 	c.HTML(http.StatusOK, "api/index.html", gin.H{
		// 		"title":     "前台首页",
		// 		"news":      a,
		// 		"hobby":     []string{"吃饭", "睡觉", "写代码"},
		// 		"date":      time.Now().Unix(),
		// 		"printInfo": "Building...",
		// 	})
		// }
		// )
	}
}
