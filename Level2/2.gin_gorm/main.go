package main

import (
	"fmt"
	"gin/models"
	"gin/routers"
	"html/template"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// 安装热加载程序：go get github.com/pilu/fresh
// 检查go和gopath/bin下有没有fresh.exe,否则找不到命令fresh
// 下载postman：https://www.postman.com/downloads/

// type UserList struct {
// 	Title   string `json:"title"`
// 	Desc    string `json:"desc"`
// 	Content string `json:"content"`
// }

// type UserInfo struct {
// 	UserName string `json:"username" form:"username" xml:"username"`
// 	Password string `json:"password" form:"password" xml:"password"`
// }

// // 时间戳转换成日期函数
//
//	func UnixToTime(timestamp int64) string {
//		fmt.Println("timestamp = ", timestamp)
//		t := time.Unix(timestamp, 0)
//		return t.Format("2006-01-02 15:04:05")
//	}
// func Println(str1 string, str2 string) string {
// 	return str1 + str2
// }

func initMiddleware(c *gin.Context) {
	fmt.Println("===这是中间件开始===")
	c.Next() // 继续后续剩余处理程序
	// c.Abort() // 终止，不会继续后续剩余处理程序
	fmt.Println("===这是中间件结束===")
}

func initMiddleware2(c *gin.Context) {
	fmt.Println("===这是中间件2开始===")
	c.Next() // 继续后续剩余处理程序
	// c.Abort() // 终止，不会继续后续剩余处理程序
	fmt.Println("===这是中间件2结束===")
}

func main() {
	// 创建一个默认的路由引擎
	router := gin.Default()
	//自定义模板函数,必须在r.LoadHTMLGlob前面
	router.SetFuncMap(template.FuncMap{
		"UnixToTime": models.UnixToTime, //注册模板函数
		"Println":    models.Println,
	})
	//加载templates中所有模板文件, 使用不同目录下名称相同的模板,注意:一定要放在配置路由之前才得行
	// 配置模板的文件
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static", "./static")

	//配置全局中间件
	//router.Use(initMiddleware,initMiddleware2)
	// 引入session中间件
	//secret123456" 加密密钥
	// store := cookie.NewStore([]byte("secret123456"))
	// 配置session中间件，store 是前面创建的存储引擎，可以替换成其他存储引擎
	// router.Use(sessions.Sessions("mysession", store))

	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// 路由配置分组
	routers.AdminRoutersInit(router)
	routers.DefaultRoutersInit(router)
	routers.ApiRoutersInit(router)
	// 配置路由
	// testTemplate(router)

	// defaultRouters := router.Group("/")
	// {
	// 	defaultRouters.GET("/", func(c *gin.Context) {
	// 		c.String(200, "%v", "你好gin")
	// 	})
	// 中间件
	router.GET("/middle", initMiddleware, initMiddleware2, func(c *gin.Context) {
		fmt.Println("===这是请求开始===")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "你好middle中间件",
		})
		fmt.Println("===这是请求结束===")
	})
	// 	defaultRouters.GET("/struct", func(c *gin.Context) {
	// 		u := &UserList{
	// 			Title:   "gin",
	// 			Desc:    "你好gin",
	// 			Content: "gin 请求",
	// 		}
	// 		c.JSON(http.StatusOK, u)
	// 	})
	// 	defaultRouters.GET("/value", func(c *gin.Context) {
	// 		//127.0.0.1:8000/value?username="张三"&age=60&page=10
	// 		username := c.Query("username")
	// 		age := c.Query("age")
	// 		page := c.DefaultQuery("page", "1")

	// 		c.JSON(http.StatusOK, gin.H{
	// 			"username": username,
	// 			"age":      age,
	// 			"page":     page,
	// 		})
	// 	})
	// 	defaultRouters.GET("/user", func(c *gin.Context) {
	// 		c.HTML(http.StatusOK, "default/user.html", gin.H{
	// 			"title": "用户首页",
	// 		})
	// 	})
	// 	defaultRouters.GET("/getUser", func(c *gin.Context) {
	// 		fmt.Println("getUser=====================")
	// 		var user UserInfo
	// 		if err := c.ShouldBind(&user); err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 			return
	// 		}
	// 		c.JSON(http.StatusOK, user)

	// 	})

	// 	defaultRouters.POST("/getUserXML", func(c *gin.Context) {
	// 		var user UserInfo
	// 		xmlData, _ := c.GetRawData()

	// 		if err := xml.Unmarshal(xmlData, &user); err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 			return
	// 		}

	// 		c.JSON(http.StatusOK, user)
	// 	})

	// 	defaultRouters.POST("/doAddUser", func(c *gin.Context) {
	// 		username := c.PostForm("username")
	// 		password := c.PostForm("password")

	// 		c.JSON(http.StatusOK, gin.H{
	// 			"username": username,
	// 			"password": password,
	// 		})
	// 	})

	// 	defaultRouters.GET("/xml", func(c *gin.Context) {
	// 		// gin.H 在 Gin 框架中经过特殊处理‌，自动补全了 XML 根元素，
	// 		// 而普通 map 因不符合 XML 规范导致渲染失败
	// 		c.XML(http.StatusOK, gin.H{
	// 			"success": true,
	// 			"message": "你好gin",
	// 		})
	// 	})

	// 	defaultRouters.POST("/add", func(c *gin.Context) {
	// 		c.String(200, "这是一个post请求")
	// 	})
	// 	defaultRouters.PUT("/edit", func(c *gin.Context) {
	// 		c.String(200, "这是一个put请求")
	// 	})
	// 	defaultRouters.DELETE("/delete", func(c *gin.Context) {
	// 		c.String(http.StatusOK, "这是一个delete请求")
	// 	})
	// }

	router.Run(":8000") // 监听并在 0.0.0.0:8080 上启动服务
}

// type Article struct {
// 	Title   string
// 	Content string
// 	Score   int
// }

// html模板
// func testTemplate(router *gin.Engine) {
// adminRouters := router.Group("/admin")
// {
// 	adminRouters.GET("/", func(c *gin.Context) {
// 		// gin.H 在 Gin 框架中经过特殊处理‌，自动补全了 XML 根元素，
// 		// 而普通 map 因不符合 XML 规范导致渲染失败
// 		c.HTML(http.StatusOK, "admin/index.html", gin.H{
// 			"title": "后台首页",
// 		})
// 	})

// 	adminRouters.GET("/news", func(c *gin.Context) {
// 		a := &model.Article{
// 			Title:   "后台新闻页面",
// 			Content: "后台新闻详情",
// 		}
// 		c.HTML(http.StatusOK, "admin/news.html", gin.H{
// 			"title": "后台新闻",
// 			"news":  a,
// 		})
// 	})
// }

// apiRouters := router.Group("/api")
// {
// 	apiRouters.GET("/news", func(c *gin.Context) {
// 		a := &model.Article{
// 			Title:   "api新闻页面",
// 			Content: "api新闻详情",
// 		}
// 		c.HTML(http.StatusOK, "api/news.html", gin.H{
// 			"title": "api新闻",
// 			"news":  a,
// 		})
// 	})

// 	apiRouters.GET("/", func(c *gin.Context) {
// 		a := &model.Article{
// 			Title:   "首页页面",
// 			Content: "首页详情",
// 			Score:   91,
// 		}
// 		c.HTML(http.StatusOK, "api/index.html", gin.H{
// 			"title":     "前台首页",
// 			"news":      a,
// 			"hobby":     []string{"吃饭", "睡觉", "写代码"},
// 			"date":      time.Now().Unix(),
// 			"printInfo": "Building...",
// 		})
// 	})
// }

// }
