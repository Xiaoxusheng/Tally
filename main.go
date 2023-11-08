package main

import (
	"Tally/config"
	"Tally/router"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func main() {
	config.InitService()

	config.InitMysql()

	config.InitRedis()

	e := echo.New()
	e.Debug = true

	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	router.Routers(e)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Config.Service.Port)))
}
