package user

import (
	"Tally/common"
	"Tally/dao"
	"Tally/global"
	"Tally/utils"
	"github.com/labstack/echo/v4"
)

/*
关注列表
*/

// FollowUser 关注用户
func FollowUser(c echo.Context) error {
	//关注人的id
	identity := c.QueryParam("identity")
	if identity == "" {
		return common.Fail(c, global.UserCode, global.QueryErr)
	}
	//用户id
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.UserCode, global.QueryErr)
	}
	//判断用户的id是否为真
	user := dao.GetUserByIdentity(id)
	if user == nil {
		return common.Fail(c, global.UserCode, global.UserNotFound)
	}

	//判断是否存在
	if global.Global.Redis.SIsMember(global.Global.Ctx, global.UserFollow+id, identity).Val() {
		return common.Fail(c, global.UserCode, global.AlreadyFollow)
	}
	//判断是否封禁
	if global.Global.Redis.SIsMember(global.Global.Ctx, global.BanUser, identity).Val() {
		return common.Fail(c, global.UserCode, global.BannedUser)
	}
	//关注
	val := global.Global.Redis.SAdd(global.Global.Ctx, global.UserFollow+id, identity).Val()
	//	加入关注列表
	if val == global.Fail {
		return common.Fail(c, global.UserCode, global.FollowFail)
	}
	return common.Ok(c, nil)
}

// CancelFollow 取消关注
func CancelFollow(c echo.Context) error {
	//关注人的id
	identity := c.QueryParam("identity")
	if identity == "" {
		return common.Fail(c, global.UserCode, global.QueryErr)
	}
	//用户id
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.UserCode, global.QueryErr)
	}
	//判断是否存在
	if global.Global.Redis.SIsMember(global.Global.Ctx, global.UserFollow, id).Val() {
		return common.Fail(c, global.UserCode, global.AlreadyFollow)
	}
	//取消关注
	val := global.Global.Redis.SRem(global.Global.Ctx, global.UserFollow, id).Val()
	//	从关注列表移除
	if val == global.Fail {
		return common.Fail(c, global.UserCode, global.CancelFollowFail)
	}
	return common.Ok(c, nil)
}