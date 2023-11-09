package user

import (
	"Tally/common"
	"Tally/config"
	"Tally/dao"
	"Tally/global"
	"Tally/utils"
	"fmt"
	"github.com/labstack/echo/v4"
)

type User struct {
	Username string `json:"username"  form:"username"  validate:"min=5,max=10"`
	Password string `json:"password"  form:"password" validate:"min=5,max=10"`
}

func Login(c echo.Context) error {
	user := new(User)
	err := c.Bind(user)
	if err != nil {
		return common.Fail(c, global.UserCode, "解析错误！")
	}
	err = c.Validate(user)
	if err != nil {
		return err
	}
	fmt.Println(user)
	//从redis获取
	val := global.Global.Redis.HGet(global.Global.Ctx, user.Username, user.Password).Val()
	token := utils.GetToken(val)
	if val != "" {
		fmt.Println(val)
		return common.Ok(c, map[string]any{
			"token": token,
		})
	} else {
		//数据库中获取
		ok := dao.GetUserById(user.Username, utils.Md5(user.Password))
		if ok == nil {
			return common.Fail(c, global.UserCode, "用户名或密码错误")
		}
		//异步更新
		go func() {
			global.Global.Redis.HSet(global.Global.Ctx, user.Username, user.Password, ok.Identity)
			//	放入token
			global.Global.Redis.Set(global.Global.Ctx, ok.Identity, token, config.Config.Jwt.Time)
		}()
		return common.Ok(c, map[string]any{
			"token": token,
		})
	}
}
