package user

import (
	"Tally/common"
	"Tally/dao"
	"Tally/global"
	"github.com/labstack/echo/v4"
)

type Searches struct {
	Identity     string  `json:"identity"  query:"identity" form:"identity"  param:"identity"`
	UserIdentity string  `json:"userIdentity" query:"userIdentity" form:"userIdentity" param:"userIdentity"`
	Kind         int     `json:"kind" query:"kind" form:"kind" param:"kind"`
	Money        float64 `json:"money" query:"money" form:"money" param:"money"`
	Remark       string  `json:"remark" query:"remark" form:"remark" param:"remark"`
	Category     int     `json:"category" query:"category" form:"category" param:"category"`
}

func Search(c echo.Context) error {
	s := new(Searches)
	err := c.Bind(s)
	if err != nil {
		return common.Fail(c, global.SearchCode, "转换失败")
	}

	if s.Remark != "" {
		list := dao.GetLikeList(s.Remark)
		return common.Ok(c, list)
	} else if s.Identity != "" {
		l := dao.GetListById(s.Identity)
		return common.Ok(c, l)
	}
	return nil
}
