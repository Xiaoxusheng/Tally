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

type Collect struct {
	Identity string `json:"identity"  query:"identity" param:"identity"`
}

// AddCollect 添加收藏
func AddCollect(c echo.Context) error {
	collect := new(Collect)
	err := c.Bind(collect)
	if err != nil {
		return common.Fail(c, global.CollectCode, global.ParseErr)
	}
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.TallyCode, global.UserIdentityErr)
	}
	if ok := dao.GetByCollect(collect.Identity, id); !ok {
		return common.Fail(c, global.TallyCode, global.CollectErr)
	}
	err = dao.UpdateCollect(collect.Identity, id)
	if err != nil {
		return common.Fail(c, global.TallyCode, global.CollectToErr)
	}
	go func() {
		global.Global.Redis.Del(global.Global.Ctx, global.CollectKey+id)
	}()
	return common.Ok(c, nil)

}

// DeleteCollect 删除收藏
func DeleteCollect(c echo.Context) error {
	collect := new(Collect)
	err := c.Bind(collect)
	if err != nil {
		return common.Fail(c, global.CollectCode, global.ParseErr)
	}
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.TallyCode, global.UserIdentityErr)
	}
	if ok := dao.GetByCollect(collect.Identity, id); !ok {
		return common.Fail(c, global.TallyCode, global.CollectErr)
	}
	err = dao.CancelCollect(collect.Identity, id)
	if err != nil {
		return common.Fail(c, global.TallyCode, global.CollectToErr)
	}
	go func() {
		global.Global.Redis.Del(global.Global.Ctx, global.CollectKey+id)
	}()
	return common.Ok(c, nil)
}

// CollectList 收藏列表
func CollectList(c echo.Context) error {
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.TallyCode, global.UserIdentityErr)
	}
	val := global.Global.Redis.Get(global.Global.Ctx, global.CollectKey+id).Val()
	if val != "" {
		tally := make([]*models.Tally, 0)
		err := json.Unmarshal([]byte(val), &tally)
		if err != nil {
			return err
		}
		return common.Ok(c, tally)
	} else {
		list := dao.GetCollectList(id)
		//异步更新

		return common.Ok(c, list)
	}
}
