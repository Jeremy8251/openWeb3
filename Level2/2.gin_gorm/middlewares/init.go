package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitMiddleware(c *gin.Context) {
	fmt.Println("===这是中间件开始===")
	c.Set("username", "zhangsan123")
	c.Next() // 继续后续剩余处理程序
	// c.Abort() // 终止，不会继续后续剩余处理程序
	fmt.Println("===这是中间件结束===")
	// cCp := c.Copy()
	// go func() {
	// 	time.Sleep(2 * time.Second)
	// 	fmt.Println("Done in path" + cCp.Request.URL.Path)
	// }()
}
