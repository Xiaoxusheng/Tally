package user

import (
	"Tally/common"
	"Tally/dao"
	"Tally/global"
	"Tally/models"
	"Tally/utils"
	"github.com/labstack/echo/v4"
	"time"
)

type List struct {
	Kind     int     `json:"kind"`
	Money    float64 `json:"money"`
	Remark   string  `json:"remark"`
	Category int     `json:"category"`
}

func TallyList(c echo.Context) error {
	id, ok := c.Get("identity").(string)
	if !ok {
		return common.Fail(c, global.TallyCode, "获取失败")
	}
	//缓存中获取
	val := global.Global.Redis.HGet(global.Global.Ctx, id, "list").Val()
	if val != "" {
		return common.Ok(c, val)
	} else {
		list := dao.GetTallyList(id)
		if list == nil {
			return common.Fail(c, global.TallyCode, "获取失败")
		}
		go func() {
			global.Global.Redis.HSet(global.Global.Ctx, id, "list", list, time.Duration(utils.GetRandom()))
		}()
		return common.Ok(c, list)
	}

}

func AddTallyLog(c echo.Context) error {
	t := new(List)
	userIdentity := utils.GetIdentity(c, "identity")
	if userIdentity == "" {
		return common.Fail(c, global.TallyCode, "获取失败")
	}
	err := c.Bind(t)
	if err != nil {
		return common.Fail(c, global.TallyCode, "参数错误")
	}
	err = dao.InsertTally(&models.Tally{
		Identity:     utils.GetUidV4(),
		UserIdentity: userIdentity,
		Kind:         t.Kind,
		Money:        t.Money,
		Remark:       t.Remark,
		Category:     t.Category,
	})
	if err != nil {
		return common.Fail(c, global.TallyCode, "添加失败")
	}
	return common.Ok(c, nil)
}

func AllotKind(c echo.Context) error {
	kind := c.QueryParam("kind")
	userIdentity := utils.GetIdentity(c, "identity")
	if userIdentity == "" {
		return common.Fail(c, global.TallyCode, "获取失败")
	}
	val := global.Global.Redis.HGet(global.Global.Ctx, userIdentity, kind).Val()
	if val != "" {
		return common.Ok(c, val)
	} else {
		list := dao.GetTallyKind(userIdentity, kind)
		go func() {
			global.Global.Redis.HSet(global.Global.Ctx, userIdentity, kind, list)
			global.Global.Redis.Expire(global.Global.Ctx, userIdentity, time.Duration(utils.GetRandom()))
		}()
		return common.Ok(c, list)
	}
}
