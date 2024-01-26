package admin

import (
	"Tally/common"
	"Tally/dao"
	"Tally/global"
	"github.com/labstack/echo/v4"
)

// GetUserInfoList 查看所有用户信息
func GetUserInfoList(c echo.Context) error {
	list, err := dao.GetUserList()
	if err != nil {
		return common.Fail(c, global.AdminCode, global.GetUserListFail)
	}
	return common.Ok(c, list)
}

// BanUser 封禁用户
func BanUser(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return common.Fail(c, global.AdminCode, global.QueryErr)
	}
	//判断id是否存在
	if dao.GetUserByIdentity(id) == nil {
		return common.Fail(c, global.AdminCode, global.UserNotfound)
	}
	//添加进封禁集合
	if global.Fail == global.Global.Redis.SAdd(global.Global.Ctx, global.BanUser, id).Val() {
		return common.Fail(c, global.AdminCode, global.BanUserFail)
	}
	global.Global.Pool.Submit(func() {
		err := dao.UpdateUserStatus(id, global.Success)
		if err != nil {
			return
		}
	})
	return common.Ok(c, nil)
}

// UnsealUser 解封用户
func UnsealUser(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return common.Fail(c, global.AdminCode, global.QueryErr)
	}
	//判断id是否存在
	if dao.GetUserByIdentity(id) == nil {
		return common.Fail(c, global.AdminCode, global.UserNotfound)
	}
	//添加进封禁集合
	if global.Fail == global.Global.Redis.SRem(global.Global.Ctx, global.BanUser, id).Val() {
		return common.Fail(c, global.AdminCode, global.BanUserFail)
	}
	global.Global.Pool.Submit(func() {
		err := dao.UpdateUserStatus(id, global.Fail)
		if err != nil {
			return
		}
	})
	return common.Ok(c, nil)
}

//下架博客

//人工审核博客

//人工审核评论

//人工审核图片

//查看用户访问量

//查看ip访问
