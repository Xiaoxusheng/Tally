package user

import (
	"Tally/common"
	"Tally/dao"
	"Tally/global"
	"Tally/models"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
)

type Kind struct {
	Name         string `json:"name,omitempty" query:"name" form:"name" param:"name" validate:"name"`
	SerialNumber int    `json:"serialNumber,omitempty" query:"serialNumber" form:"serialNumber" param:"serialNumber"  validate:"required"`
}

// AddKind 添加种类
func AddKind(c echo.Context) error {
	k := new(Kind)
	err := c.Bind(k)
	if err != nil {
		return common.Fail(c, global.KindCode, global.ParseErr)
	}
	err = dao.InsertKind(&models.Kind{
		Name:         k.Name,
		SerialNumber: k.SerialNumber,
	})
	if err != nil {
		return common.Fail(c, global.KindCode, "添加失败")
	}
	//删除数据
	go func() {
		global.Global.Redis.Del(global.Global.Ctx, global.ListKey)
	}()

	return common.Ok(c, nil)
}

// KindList 种类列表
func KindList(c echo.Context) error {
	val := global.Global.Redis.Get(global.Global.Ctx, global.ListKey).Val()
	if val != "" {
		fmt.Println("val不为空")
		var s []models.Kind
		err := json.Unmarshal([]byte(val), &s)
		if err != nil {
			return common.Fail(c, global.KindCode, "数据错误")
		}
		return common.Ok(c, s)
	} else {
		list := dao.GetKindList()
		go func() {
			marshal, err := json.Marshal(list)
			if err != nil {
				global.Global.Log.Warn("序列化失败！")
			}
			result, err := global.Global.Redis.Set(global.Global.Ctx, global.ListKey, marshal, 0).Result()
			if err != nil {
				return
			}
			fmt.Println(result, err)
		}()
		return common.Ok(c, list)
	}

}
