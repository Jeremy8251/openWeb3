package admin

import (
	"fmt"
	"gin/models"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
}

func (con UserController) Index(c *gin.Context) {

	con.success(c)
	//con.error(c)
	// c.HTML(http.StatusOK, "admin/index.html", gin.H{
	// 	"title": "后台首页",
	// })
}

func (con UserController) Job(c *gin.Context) {

	jobinfoList := []models.Jobinfo{}
	// 查询第一条数据
	// models.DB.First(&jobinfoList)

	models.DB.Where("jobaddress like ?", "%工作地点：澳门区%").Find(&jobinfoList)
	c.JSON(http.StatusOK, gin.H{
		"result": jobinfoList,
	})
}

func (con UserController) News(c *gin.Context) {
	a := &models.Article{
		Title:   "后台新闻页面",
		Content: "后台新闻详情",
	}
	c.HTML(http.StatusOK, "admin/news.html", gin.H{
		"title": "后台新闻",
		"news":  a,
	})

}

func (con UserController) Add(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/useradd.html", nil)

}

func (con UserController) DoUpLoad(c *gin.Context) {
	username := c.PostForm("username")
	file, err := c.FormFile("face")

	dst := path.Join("./static/upload", file.Filename)
	fmt.Println("dst=======", dst)
	if err == nil {
		c.SaveUploadedFile(file, dst)
	}
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"username": username,
		"dst":      dst,
	})

}

func (con UserController) DoEdit(c *gin.Context) {
	username := c.PostForm("username")
	face1, err1 := c.FormFile("face1")
	face2, err2 := c.FormFile("face2")

	dst1 := path.Join("./static/upload", face1.Filename)
	if err1 == nil {
		c.SaveUploadedFile(face1, dst1)
	}

	dst2 := path.Join("./static/upload", face2.Filename)
	if err2 == nil {
		c.SaveUploadedFile(face2, dst2)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"username": username,
		"dst1":     dst1,
		"dst2":     dst2,
	})

}
