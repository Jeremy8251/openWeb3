package routers

import (
	"gin/controllers/admin"
	"gin/middlewares"

	"github.com/gin-gonic/gin"
)

func AdminRoutersInit(router *gin.Engine) {
	// 配置这个中间件只在这个路由生效
	adminRouters := router.Group("/admin", middlewares.InitMiddleware)
	{
		adminRouters.GET("/", admin.UserController{}.Index)

		adminRouters.GET("/job", admin.UserController{}.Job)
		// func(c *gin.Context) {
		// 	// gin.H 在 Gin 框架中经过特殊处理‌，自动补全了 XML 根元素，
		// 	// 而普通 map 因不符合 XML 规范导致渲染失败
		// 	c.HTML(http.StatusOK, "admin/index.html", gin.H{
		// 		"title": "后台首页",
		// 	})
		// })

		adminRouters.GET("/news", admin.UserController{}.News)
		// func(c *gin.Context) {

		// a := &model.Article{
		// 	Title:   "后台新闻页面",
		// 	Content: "后台新闻详情",
		// }
		// c.HTML(http.StatusOK, "admin/news.html", gin.H{
		// 	"title": "后台新闻",
		// 	"news":  a,
		// })
		// }
		// )

		adminRouters.GET("/user/add", admin.UserController{}.Add)

		adminRouters.POST("/user/doUpload", admin.UserController{}.DoUpLoad)

		adminRouters.POST("/user/doEdit", admin.UserController{}.DoEdit)
	}
}
