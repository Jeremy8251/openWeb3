package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var err error

func init() {
	// DSN:Data Source Name
	dsn := "root:123456@tcp(127.0.0.1:3306)/jobinfo?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound错误
			Colorful:                  false,       // 禁用彩色输出
		},
	)
	// 打开数据库连接（但不校验）
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// NamingStrategy: schema.NamingStrategy{
		// 	SingularTable: true, // 使用单数表名
		// },
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println("sql.Open failed", err)
	}

}
