package user

import (
	"Tally/common"
	"Tally/dao"
	"Tally/global"
	"Tally/models"
	"Tally/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
)

type List struct {
	Kind     int     `json:"kind" query:"kind" form:"kind"`
	Money    float64 `json:"money" query:"money" form:"money"`
	Remark   string  `json:"remark" query:"remark" form:"remark"`
	Category int     `json:"category" query:"category" form:"category"`
}

type Time struct {
	StarTime int64 `query:"star_time"`
	EndTime  int64 `query:"end_time"`
}

// TallyList 获取所有账单
func TallyList(c echo.Context) error {
	id, ok := c.Get("identity").(string)
	if !ok {
		return common.Fail(c, global.TallyCode, "获取失败")
	}
	//缓存中获取
	val := global.Global.Redis.Get(global.Global.Ctx, id+"list").Val()
	fmt.Println("val", val)
	if val != "" {
		var tally []models.Tally
		err := json.Unmarshal([]byte(val), &tally)
		if err != nil {
			return err
		}
		return common.Ok(c, tally)
	} else {
		list := dao.GetTallyList(id)
		if list == nil {
			return common.Fail(c, global.TallyCode, "获取失败")
		}
		go func() {
			marshal, err := json.Marshal(list)
			if err != nil {
				return
			}
			val, err := global.Global.Redis.Set(global.Global.Ctx, id+"list", marshal, 0).Result()
			fmt.Println(val, err)
		}()
		return common.Ok(c, list)
	}

}

// AddTallyLog 添加账单
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
	go func() {
		val, err := global.Global.Redis.Del(global.Global.Ctx, userIdentity+"list").Result()
		fmt.Println("删除", val, err)
	}()
	return common.Ok(c, nil)
}

// AllotKind 根据分类获取列表数据
func AllotKind(c echo.Context) error {
	kind := c.QueryParam("kind")
	userIdentity := utils.GetIdentity(c, "identity")
	if userIdentity == "" {
		return common.Fail(c, global.TallyCode, "获取失败")
	}
	val := global.Global.Redis.Get(global.Global.Ctx, userIdentity+kind).Val()
	if val != "" {
		var list []models.Tally
		err := json.Unmarshal([]byte(val), &list)
		if err != nil {
			return err
		}
		return common.Ok(c, list)
	} else {
		list := dao.GetTallyKind(userIdentity, kind)
		go func() {
			marshal, err := json.Marshal(list)
			if err != nil {
				return
			}
			//存入redis
			global.Global.Redis.Set(global.Global.Ctx, userIdentity+kind, marshal, 0)
		}()
		return common.Ok(c, list)
	}
}

// DateList 按日期查询
func DateList(c echo.Context) error {
	//传时间戳
	t := new(Time)
	userIdentity := utils.GetIdentity(c, "identity")
	if userIdentity == "" {
		return common.Fail(c, global.TallyCode, "获取失败")
	}
	err := c.Bind(t)
	if err != nil {
		return common.Fail(c, global.TallyCode, "参数错误")
	}
	if t.StarTime > t.EndTime {
		return common.Fail(c, global.TallyCode, "时间错误")
	}
	star := time.Unix(t.StarTime, 0)
	end := time.Unix(t.EndTime, 0)
	fmt.Println(star, end)
	//redis
	val := global.Global.Redis.Get(global.Global.Ctx, userIdentity+star.String()+end.String()).Val()
	if val != "" {
		var vals []models.Tally
		err := json.Unmarshal([]byte(val), &vals)
		if err != nil {
			return err
		}
		return common.Ok(c, vals)
	} else {
		list := dao.GetByTime(star.String(), end.String())
		if list == nil {
			//没查到，防止穿透
			go func() {
				global.Global.Redis.Set(global.Global.Ctx, userIdentity+star.String()+end.String(), "null", 0)
			}()
			return common.Fail(c, global.TallyCode, "查询失败")
		}
		go func() {
			marshal, err := json.Marshal(list)
			if err != nil {
				return
			}
			global.Global.Redis.Set(global.Global.Ctx, userIdentity+star.String()+end.String(), marshal, 0)
		}()
		return common.Ok(c, list)
	}
}

// BindKind 绑定分类
func BindKind(c echo.Context) error {
	kind, err := strconv.Atoi(c.QueryParam("kind "))
	if err != nil {
		return common.Fail(c, global.TallyCode, "转换失败")
	}
	id := c.QueryParam("identity ")
	useIdentity := utils.GetIdentity(c, "identity")
	if useIdentity == "" {
		return common.Fail(c, global.TallyCode, "获取失败")
	}
	if ok := dao.GetByKind(kind); !ok {
		return common.Fail(c, global.TallyCode, "分类不存在")
	}
	err = dao.UpdateByKind(id, kind)
	if err != nil {
		return common.Fail(c, global.TallyCode, "绑定失败")
	}

	return common.Ok(c, nil)
}
