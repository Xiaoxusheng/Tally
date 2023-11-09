package user

import (
	"Tally/common"
	"Tally/dao"
	"Tally/global"
	"Tally/models"
	"Tally/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type RegisterUser struct {
	Username string `json:"username"  form:"username"  validate:"min=5,max=10"`
	Password string `json:"password"  form:"password" validate:"min=5,max=10"`
	Phone    string `json:"phone" form:"phone" validate:"required"`
}

func Register(c echo.Context) error {
	user := new(RegisterUser)
	err := c.Bind(user)
	if err != nil {
		return common.Fail(c, http.StatusOK, "解析失败！")
	}
	err = c.Validate(user)
	if err != nil {
		return err
	}
	//注册
	if ok := dao.GetPhone(user.Phone); ok {
		return common.Fail(c, http.StatusOK, "电话号码已经存在！")
	}
	//
	if ok := dao.GetUserByUsername(user.Username); ok {
		return common.Fail(c, http.StatusOK, "用户名已经存在！")
	}
	id := utils.GetUidV5(user.Username)
	err = dao.InsertUser(&models.User{
		Username: user.Username,
		Phone:    user.Phone,
		Password: utils.Md5(user.Password),
		IP:       c.RealIP(),
		Identity: id,
	})
	if err != nil {
		return common.Fail(c, http.StatusOK, "注册失败！")
	}
	go func() {
		global.Global.Redis.HSet(global.Global.Ctx, user.Username, user.Password, id)
	}()
	return common.Ok(c, nil)
}
