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
	e.Use(middleware.Logger(), middleware.CORS(), middleware.Timeout())
	//middleware.Recover()

	e.POST("/user/login", user.Login)
	e.POST("/user/register", user.Register)
	e.GET("/oauth/redirect", user.Token)

	users := e.Group("/user")
	users.GET("/auth2_login", user.OAuthLogin)
	users.Use(m.ParseToken())
	users.GET("/list", user.TallyList)
	users.GET("/info", user.Info)
	users.GET("/change_pwd", user.ChangePwd)
	users.GET("/logout", user.Logout)
	users.POST("/add", user.AddTallyLog)
	users.GET("/allot_kind", user.AllotKind)
	users.GET("/bind", user.BindKind)
	users.GET("/date_list", user.DateList)
	users.Any("/addKind", user.AddKind)
	users.GET("/kind_list", user.KindList)
	users.GET("/search", user.Search)
	users.GET("/add_collect", user.AddCollect)
	users.GET("/del_collect", user.DeleteCollect)
	users.GET("/collect_list", user.CollectList)
	users.GET("/analysis", user.Analysis)
	users.POST("/upload", user.Upload)
	users.POST("/blog_text", user.BlogText)
	users.GET("/likes", user.Likes)

}
