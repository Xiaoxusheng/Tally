package user

import (
	"Tally/common"
	"Tally/config"
	"Tally/dao"
	"Tally/global"
	"Tally/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"time"
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
	val := global.Global.Redis.HGet(global.Global.Ctx, user.Username, utils.Md5(user.Password)).Val()
	token := utils.GetToken(val)
	if val != "" {
		fmt.Println("identity得值", val)
		go func() {
			global.Global.Redis.Set(global.Global.Ctx, val, token, config.Config.Jwt.Time*time.Hour)
		}()
		return common.Ok(c, map[string]any{
			"token": token,
		})
	} else {
		//数据库中获取
		ok := dao.GetUserById(user.Username, utils.Md5(user.Password))
		if ok == nil {
			return common.Fail(c, global.UserCode, "用户名或密码错误")
		}
		token := utils.GetToken(ok.Identity)
		//异步更新
		go func() {
			global.Global.Redis.HSet(global.Global.Ctx, user.Username, user.Password, ok.Identity)
			//	放入token
			global.Global.Redis.Set(global.Global.Ctx, ok.Identity, token, config.Config.Jwt.Time*time.Hour)
		}()
		return common.Ok(c, map[string]any{
			"token": token,
		})
	}
}

func Logout(c echo.Context) error {
	identity := utils.GetIdentity(c, "identity")
	if identity == "" {
		return common.Fail(c, global.UserCode, "获取失败")
	}
	//删除token
	val := global.Global.Redis.Del(global.Global.Ctx, identity).Val()
	if val != 1 {
		return common.Fail(c, global.UserCode, "退出失败！")
	}
	return common.Ok(c, nil)
}

func Info(c echo.Context) error {
	identity := utils.GetIdentity(c, "identity")
	if identity == "" {
		return common.Fail(c, global.UserCode, "获取失败")
	}
	//redis获取
	val := global.Global.Redis.Get(global.Global.Ctx, identity+"info").Val()
	fmt.Println("redis获取", val)
	if val != "" {
		return common.Ok(c, val)
	} else {
		info := dao.GetUserByIdentity(identity)
		if info == nil {
			return common.Fail(c, global.UserCode, "获取个人信息失败")
		}
		go func() {
			val, err := global.Global.Redis.Set(global.Global.Ctx, identity+"info", info, 0).Result()
			fmt.Println(val, err)
		}()
		return common.Ok(c, info)
	}
}
