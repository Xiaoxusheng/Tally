package router

import (
	"Tally/config"
	"Tally/controller/user"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routers(e *echo.Echo) {
	e.Use(middleware.Logger(), middleware.CORS(), middleware.Timeout(), middleware.Recover())

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.Config.Jwt.Key),
	}))

	users := e.Group("/user")
	users.GET("/login", user.Login)
}
