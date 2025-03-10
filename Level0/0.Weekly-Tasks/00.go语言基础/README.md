### 来源github
# Golang 语言重点学习路线

# Go 语言环境搭建

## 环境变量

GOROOT

GOPATH

GOPROXY

export PATH=$GOROOT/bin:$GOPATH/bin:$PATH

## 包管理

### Mod

添加依赖

go get <package>

更新依赖

go get -u <package>

整理依赖

go mod tidy

替换依赖为其他包

go mod edit -replace < 原 package>=< 新 package 或本地路径 >

# 语言基础重点

## 函数与闭包

## 切片

append 关键字

copy 关键字

range 关键字

## goroutine 使用

启动

结束

## channel 使用

声明

传递

读取与写入

配合 select 关键字

## defer 关键字

搭配 sync.Mutex 使用

## context 包

声明

context.WithCancel

context.WithTimeout

## sync 包工具

sync.Mutex

sync.Once

sync.Pool

sync.WaitGroup

# 框架学习

## gin

官方文档 [https://gin-gonic.com/docs/examples/](https://gin-gonic.com/docs/examples/)

配置路由(声明公开 API)

配置中间件(类似 tomcat 的 filter 和 spring 中的 interceptor)

配置日志

读取配置

## gorm

官方文档 [https://gorm.io/docs](https://gorm.io/docs)

连接数据库

基本 crud

执行事务

## ethclient(sdk 使用)

查询区块

查询交易

查询 receipt

执行交易，部署合约，调用合约（非生成代码，手动构造 input 参数，知道合约调用的本质）

查询余额

查询 account

查询 event，过滤合约事件

# CI/CD 基础操作

## Makefile

go test

go build

check lint

build docker

Dockerfile

Github action 使用
