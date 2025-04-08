package models

import "pledge-backend/db"

// 自动快速同步模型和数据库表结构
func InitTable() {
	db.Mysql.AutoMigrate(&MultiSign{})
	db.Mysql.AutoMigrate(&TokenInfo{})
	db.Mysql.AutoMigrate(&TokenList{})
	db.Mysql.AutoMigrate(&PoolData{})
	db.Mysql.AutoMigrate(&PoolBases{})
}
