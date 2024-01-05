package main

import (
	"Tally/config"
	"Tally/router"
	"Tally/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
)

// 数据预热
func init() {

}

func main() {
	//读取配置文件
	config.InitService()
	//log初始化
	config.InitLog()
	//连接mysql
	config.InitMysql()
	//连接redis
	config.InitRedis()
	//初始化协程池
	config.InitPool()

	//异步写入数据库
	go utils.Set()

	//监听系统信号
	go utils.Listen()

	e := echo.New()
	e.Debug = true

	e.Use(middleware.Logger(), middleware.Recover(), middleware.CORS(), middleware.RequestID())
	router.Routers(e)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Config.Service.Port)))
}
