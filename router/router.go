package router

import (
	"Tally/controller/user"
	m "Tally/middleware"
	"Tally/utils"
	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routers(e *echo.Echo) {
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	e.Use(middleware.Logger(), middleware.CORS(), middleware.Timeout(), middleware.Recover())

	users := e.Group("/user")
	users.Use(m.ParseToken())
	e.POST("/login", user.Login)
	e.POST("/register", user.Register)
	users.GET("/list", user.TallyList)
	users.POST("/add", user.AddTallyLog)
}