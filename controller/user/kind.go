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
	Name         string `json:"name,omitempty" query:"name" form:"name" param:"name"`
	SerialNumber int    `json:"serialNumber,omitempty" query:"serialNumber" form:"serialNumber" param:"serialNumber"`
}

func AddKind(c echo.Context) error {
	k := new(Kind)
	err := c.Bind(k)
	if err != nil {
		return common.Fail(c, global.KindCode, "参数错误")
	}
	err = dao.InsertKind(&models.Kind{
		Name:         k.Name,
		SerialNumber: k.SerialNumber,
	})
	if err != nil {
		return common.Fail(c, global.KindCode, "添加失败")
	}
	return common.Ok(c, nil)
}

func KindList(c echo.Context) error {
	val := global.Global.Redis.Get(global.Global.Ctx, "kind_list").Val()
	if val != "" {
		var s []models.Tally
		err := json.Unmarshal([]byte(val), &s)
		if err != nil {
			return common.Fail(c, global.KindCode, "数据错误")
		}
		return common.Ok(c, s)
	} else {
		list := dao.GetKindList()
		fmt.Println("list", list)
		go func() {
			global.Global.Redis.Set(global.Global.Ctx, "kind_list", list, 0)
		}()
		return common.Ok(c, list)
	}

}
