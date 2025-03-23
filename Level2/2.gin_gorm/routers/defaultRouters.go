package routers

import (
	"encoding/xml"
	"fmt"
	"gin/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func DefaultRoutersInit(router *gin.Engine) {
	defaultRouters := router.Group("/")
	{
		defaultRouters.GET("/", func(c *gin.Context) {
			// 设置cookie
			c.SetCookie("cookieName", "我是cookie", 360, "/", c.Request.Host, false, false)
			// 删除cookie，maxAge置为 -1
			// c.SetCookie("cookieName", "我是cookie", -1, "/", c.Request.Host, false, false)
			// 设置sessions
			session := sessions.Default(c)
			session.Set("username", "张三666")
			session.Save() //必须调用生效

			c.String(200, "%v", "你好gin")

		})
		defaultRouters.GET("/test", func(c *gin.Context) {
			session := sessions.Default(c)
			username := session.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "你好gin",
				"session": username,
			})
		})
		defaultRouters.GET("/struct", func(c *gin.Context) {
			u := &models.UserList{
				Title:   "gin",
				Desc:    "你好gin",
				Content: "gin 请求",
			}
			c.JSON(http.StatusOK, u)
		})
		defaultRouters.GET("/value", func(c *gin.Context) {
			//127.0.0.1:8000/value?username="张三"&age=60&page=10
			username := c.Query("username")
			age := c.Query("age")
			page := c.DefaultQuery("page", "1")

			c.JSON(http.StatusOK, gin.H{
				"username": username,
				"age":      age,
				"page":     page,
			})
		})
		defaultRouters.GET("/user", func(c *gin.Context) {
			c.HTML(http.StatusOK, "default/user.html", gin.H{
				"title": "用户首页",
			})
		})
		defaultRouters.GET("/getUser", func(c *gin.Context) {
			fmt.Println("getUser=====================")
			var user models.UserInfo
			if err := c.ShouldBind(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, user)

		})

		defaultRouters.POST("/getUserXML", func(c *gin.Context) {
			var user models.UserInfo
			xmlData, _ := c.GetRawData()

			if err := xml.Unmarshal(xmlData, &user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, user)
		})

		defaultRouters.POST("/doAddUser", func(c *gin.Context) {
			username := c.PostForm("username")
			password := c.PostForm("password")

			c.JSON(http.StatusOK, gin.H{
				"username": username,
				"password": password,
			})
		})

		defaultRouters.GET("/xml", func(c *gin.Context) {
			// gin.H 在 Gin 框架中经过特殊处理‌，自动补全了 XML 根元素，
			// 而普通 map 因不符合 XML 规范导致渲染失败
			c.XML(http.StatusOK, gin.H{
				"success": true,
				"message": "你好gin",
			})
		})

		defaultRouters.POST("/add", func(c *gin.Context) {
			c.String(200, "这是一个post请求")
		})
		defaultRouters.PUT("/edit", func(c *gin.Context) {
			c.String(200, "这是一个put请求")
		})
		defaultRouters.DELETE("/delete", func(c *gin.Context) {
			c.String(http.StatusOK, "这是一个delete请求")
		})
	}
}
