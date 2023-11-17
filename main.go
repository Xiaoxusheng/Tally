package main

import (
	"Tally/config"
	"Tally/global"
	"Tally/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
)

func main() {
	//读取配置文件
	config.InitService()
	//连接mysql
	config.InitMysql()
	//连接redis
	config.InitRedis()
	//log初始化
	config.InitLog()
	global.Global.Log.Warn("服务启动成功")

	e := echo.New()
	e.Debug = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	router.Routers(e)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Config.Service.Port)))
}
