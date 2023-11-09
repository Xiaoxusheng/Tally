package router

import (
	"Tally/controller/user"
	"Tally/utils"
	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routers(e *echo.Echo) {
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	e.Use(middleware.Logger(), middleware.CORS(), middleware.Timeout(), middleware.Recover())

	users := e.Group("/user")
	users.POST("/login", user.Login)
	users.POST("/register", user.Register)
}
