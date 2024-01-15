package user

import (
	"Tally/common"
	"Tally/config"
	"Tally/dao"
	"Tally/global"
	"Tally/models"
	"Tally/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
用户模块
*/

type RegisterUser struct {
	Username string `json:"username"  form:"username"  validate:"min=5,max=10"`
	Password string `json:"password"  form:"password" validate:"min=5,max=10"`
	Phone    string `json:"phone" form:"phone" validate:"required,phone"`
}

type User struct {
	Username string `json:"username"  form:"username"  validate:"min=5,max=10"`
	Password string `json:"password"  form:"password" validate:"min=5,max=10"`
}

// Register 注册
func Register(c echo.Context) error {
	user := new(RegisterUser)
	err := c.Bind(user)
	if err != nil {
		return common.Fail(c, global.UserCode, global.ParseErr)
	}
	err = c.Validate(user)
	if err != nil {
		return err
	}
	//注册
	if ok := dao.GetPhone(user.Phone); ok {
		return common.Fail(c, global.UserCode, "电话号码已经存在！")
	}
	//
	if ok := dao.GetUserByUsername(user.Username); ok {
		return common.Fail(c, http.StatusOK, "用户名已经存在！")
	}
	id := utils.GetUidV5(user.Username)
	err = dao.InsertUser(&models.User{
		Username: user.Username,
		Account:  utils.GetAccount(),
		Password: utils.Md5(user.Password),
		Phone:    user.Phone,
		Identity: id,
		GithubId: " ",
		Status:   0,
		IsHide:   false,
		IP:       c.RealIP(),
	})
	if err != nil {
		global.Global.Log.Warn(err)
		return common.Fail(c, global.UserCode, "注册失败！")
	}
	global.Global.Pool.Submit(func() {
		global.Global.Redis.HSet(global.Global.Ctx, user.Username, user.Password, id)
	})
	return common.Ok(c, nil)
}

// SignOut 注销
func SignOut(c echo.Context) error {
	identity := utils.GetIdentity(c, "identity")
	if identity == "" {
		return common.Fail(c, global.UserCode, global.UserIdentityErr)
	}
	err := dao.DeleteUser(identity)
	if err != nil {
		return common.Fail(c, global.UserCode, global.UserIdentityErr)
	}
	//删除数据
	global.Global.Pool.Submit(func() {
		/*
			删除redis
		*/
		//删除关注
		global.Global.Redis.Del(global.Global.Ctx, global.UserFollow+identity)
		//删除个人信息
		global.Global.Redis.Get(global.Global.Ctx, identity+"info")
		//删除token
		global.Global.Redis.Del(global.Global.Ctx, identity)

		/*
			删除mysql数据
		*/

		err = dao.DeleteTally(identity)
		if err != nil {
			global.Global.Log.Warn(err)
		}
		err = dao.DeleteBlogCollect(identity)
		if err != nil {
			global.Global.Log.Warn(err)
		}
		err = dao.DeleteBlogCollectByUserIdentity(identity)
		if err != nil {
			global.Global.Log.Warn(err)
		}
	})
	return common.Ok(c, nil)
}

// Login 登录
func Login(c echo.Context) error {
	user := new(User)
	err := c.Bind(user)
	if err != nil {
		return common.Fail(c, global.UserCode, global.ParseErr)
	}
	err = c.Validate(user)
	if err != nil {
		return err
	}
	fmt.Println(user)
	//从redis获取
	val := global.Global.Redis.HGet(global.Global.Ctx, user.Username, utils.Md5(user.Password)).Val()
	fmt.Println("val", val)
	//获取token
	res := global.Global.Redis.Get(global.Global.Ctx, val).Val()
	if res != "" {
		// 打卡
		global.Global.Redis.SetBit(global.Global.Ctx, global.SignIn+val, int64(time.Now().Day()-1), 1)
		return common.Ok(c, map[string]any{"token": res})
	} else {
		token := utils.GetToken(val)
		if val != "" {
			global.Global.Log.Info("identity的值", val)
			global.Global.Pool.Submit(func() {
				global.Global.Redis.Set(global.Global.Ctx, val, token, config.Config.Jwt.Time*time.Hour)
				//打卡
				result, err := global.Global.Redis.SetBit(global.Global.Ctx, global.SignIn+val, int64(time.Now().Day()-1), 1).Result()
				if err != nil {
					global.Global.Log.Warn(result, err)
					return
				}
			})
			return common.Ok(c, map[string]any{"token": token})
		} else {
			//数据库中获取
			ok := dao.GetUserById(user.Username, utils.Md5(user.Password))
			if ok == nil {
				return common.Fail(c, global.UserCode, global.LoginErr)
			}
			token = utils.GetToken(ok.Identity)
			//异步更新
			global.Global.Pool.Submit(func() {
				global.Global.Redis.HSet(global.Global.Ctx, user.Username, utils.Md5(user.Password), ok.Identity)
				//	放入token
				global.Global.Redis.Set(global.Global.Ctx, ok.Identity, token, config.Config.Jwt.Time*time.Hour)
				//	打卡
				global.Global.Redis.SetBit(global.Global.Ctx, global.SignIn+ok.Identity, int64(time.Now().Day()-1), 1)
			})
			return common.Ok(c, map[string]any{
				"token": token,
			})
		}
	}
}

// Logout 登出
func Logout(c echo.Context) error {
	identity := utils.GetIdentity(c, "identity")
	if identity == "" {
		return common.Fail(c, global.UserCode, global.UserIdentityErr)
	}
	//删除token
	val := global.Global.Redis.Del(global.Global.Ctx, identity).Val()
	if val != 1 {
		return common.Fail(c, global.UserCode, "退出失败！")
	}
	return common.Ok(c, nil)
}

// Info 获取用户信息
func Info(c echo.Context) error {
	identity := utils.GetIdentity(c, "identity")
	if identity == "" {
		return common.Fail(c, global.UserCode, global.UserIdentityErr)
	}
	//redis获取
	val := global.Global.Redis.Get(global.Global.Ctx, identity+"info").Val()
	if val != "" {
		user := new(models.User)
		err := json.Unmarshal([]byte(val), user)
		if err != nil {
			return err
		}
		return common.Ok(c, user)
	} else {
		info := dao.GetUserByIdentity(identity)
		if info == nil {
			return common.Fail(c, global.UserCode, "获取个人信息失败")
		}
		global.Global.Pool.Submit(func() {
			marshal, err := json.Marshal(info)
			if err != nil {
				global.Global.Log.Warn(err)
				return
			}
			val, err := global.Global.Redis.Set(global.Global.Ctx, identity+"info", marshal, global.InfoTime).Result()
			global.Global.Log.Info(val, err)
		})
		return common.Ok(c, info)
	}
}

// ChangePwd 修改密码
func ChangePwd(c echo.Context) error {
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		global.Global.Log.Warn("identity is null")
		return common.Fail(c, global.UserCode, global.QueryErr)
	}
	OldPwd := c.QueryParam("oldpwd")
	pwd := c.QueryParam("pwd")
	if pwd == "" {
		global.Global.Log.Warn("pwd is null")
		return common.Fail(c, global.UserCode, global.PassISNull)
	}
	ok := dao.GetByPwdIdentity(id, utils.Md5(OldPwd))
	if ok == nil {
		return common.Fail(c, global.UserCode, global.PasswordIeErr)
	}
	err := dao.UpdatePwd(id, utils.Md5(pwd))
	if err != nil {
		return common.Fail(c, global.UserCode, global.ChangePassword)
	}
	//删除redis中存放的信息
	global.Global.Pool.Submit(func() {
		global.Global.Redis.Del(global.Global.Ctx, ok.Username)
	})
	//删除缓存
	return common.Ok(c, nil)
}

// OAuthLogin OAuth2 认证
func OAuthLogin(c echo.Context) error {
	configs := utils.NewOauth2()
	url := configs.AuthCodeURL(utils.GetLetter(), oauth2.SetAuthURLParam("grant_type", "Authorization Code Grant"))
	return common.Ok(c, url)
}

func Token(c echo.Context) error {
	code := c.QueryParam("code")
	global.Global.Log.Info(code)
	configs := utils.NewOauth2()
	token, err := configs.Exchange(global.Global.Ctx, code, oauth2.SetAuthURLParam("grant_type", "Authorization Code Grant"))
	if err != nil {
		global.Global.Log.Warn(err)
		return common.Fail(c, global.UserCode, "获取token失败")
	}
	//注册新用户

	//换取本地token
	//utils.GetToken()

	global.Global.Log.Info(token)
	//global.Global.Redis.Set(global.Global.Ctx, token.AccessToken, token.RefreshToken)
	return common.Ok(c, token)

	//client := oauth2.NewClient(global.Global.Ctx, oauth2.StaticTokenSource(token))
	//get, err := client.Get("https://api.github.com/user")
	//if err != nil {
	//	global.Global.Log.Warn(err)
	//	return err
	//}
	//defer func(Body io.ReadCloser) {
	//	err := Body.Close()
	//	if err != nil {
	//		global.Global.Log.Warn(err)
	//	}
	//}(get.Body)
	//all, err := io.ReadAll(get.Body)
	//if err != nil {
	//	global.Global.Log.Warn(err)
	//	return err
	//}

}

// ChangeUserInfo 修改用户信息
func ChangeUserInfo(c echo.Context) error {
	//获取用户信息
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		global.Global.Log.Warn("identity is null")
		return common.Fail(c, global.UserCode, global.QueryErr)
	}
	user := new(global.User)
	err := c.Bind(user)
	if err != nil {
		global.Global.Log.Warn(err)
		return common.Fail(c, global.UserCode, global.ParseErr)
	}
	err = dao.UpdateAll(id, user)
	if err != nil {
		global.Global.Log.Warn(err)
		return common.Fail(c, global.UserCode, global.ChangeUserInfo)
	}
	global.Global.Pool.Submit(func() {
		_, err = global.Global.Redis.Del(global.Global.Ctx, id+"info").Result()
		if err != nil {
			global.Global.Log.Warn(err)
		}
	})
	return common.Ok(c, nil)
}

// LoginInfo 获取登陆情况
func LoginInfo(c echo.Context) error {
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		global.Global.Log.Warn("identity is null")
		return common.Fail(c, global.UserCode, global.QueryErr)
	}
	d := time.Now().Day()
	val := global.Global.Redis.BitField(global.Global.Ctx, global.SignIn+id, "GET", "u"+strconv.Itoa(d), 0).Val()
	global.Global.Log.Info(val)
	s := strings.Builder{}
	if val[0] == 1 {
		for i := 1; i < d; i++ {
			s.WriteString("0")
		}
		s.WriteString("1")
		return common.Ok(c, s.String())
	}
	return common.Ok(c, fmt.Sprintf("%b", val[0]))

}
