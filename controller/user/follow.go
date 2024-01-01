package user

import (
	"Tally/common"
	"Tally/dao"
	"Tally/global"
	"Tally/models"
	"Tally/utils"
	"encoding/json"
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
	//判断是否自己
	if id == identity {
		return common.Fail(c, global.UserCode, global.FollowNot)
	}
	//判断用户的id是否为真
	if !global.Global.Redis.SIsMember(global.Global.Ctx, global.UserFollow, identity).Val() {
		user := dao.GetUserByIdentity(identity)
		global.Global.Log.Warn("user", user)
		if user == nil {
			return common.Fail(c, global.UserCode, global.UserNotFound)
		}
		val := global.Global.Redis.SAdd(global.Global.Ctx, global.UserFollow, identity).Val()
		if val == global.Fail {
			global.Global.Log.Warn("添加进redis失败")
		}
	}
	//判断是否已经关注
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
	//判断是否自己
	if id == identity {
		return common.Fail(c, global.UserCode, global.FollowNot)
	}
	//判断是否关注
	if !global.Global.Redis.SIsMember(global.Global.Ctx, global.UserFollow+id, identity).Val() {
		return common.Fail(c, global.UserCode, global.AlreadyCancelFollow)
	}
	//取消关注,从关注列表移除
	val := global.Global.Redis.SRem(global.Global.Ctx, global.UserFollow+id, identity).Val()
	if val == global.Fail {
		return common.Fail(c, global.UserCode, global.CancelFollowFail)
	}
	return common.Ok(c, nil)
}

// GetFollowList 关注列表
func GetFollowList(c echo.Context) error {
	//用户id
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.UserCode, global.QueryErr)
	}
	//获取关注用户的长度
	list := global.Global.Redis.SMembers(global.Global.Ctx, global.UserFollow+id).Val()
	global.Global.Log.Warn(list)
	if len(list) > 0 {
		followList := make([]*models.User, 0, len(list))
		for i := 0; i < len(list); i++ {
			val := global.Global.Redis.Get(global.Global.Ctx, list[i]+"info").Val()
			if val != "" {
				user := new(models.User)
				err := json.Unmarshal([]byte(val), user)
				if err != nil {
					return err
				}
				followList = append(followList, user)
			}
		}
		return common.Ok(c, followList)
	}
	followList := dao.GetFollowList(id)
	return common.Ok(c, followList)
}
