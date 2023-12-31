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
	"io"
	"strconv"
	"strings"
	"time"
)

/*
账单模块
*/

type List struct {
	Kind     int     `json:"kind" query:"kind" form:"kind"`
	Money    float64 `json:"money" query:"money" form:"money"`
	Remark   string  `json:"remark" query:"remark" form:"remark"`
	Category int     `json:"category" query:"category" form:"category"`
}

type Time struct {
	StarTime int64 `query:"star_time" json:"star_time" query:"star_time"  form:"star_time"`
	EndTime  int64 `query:"end_time" json:"end_time" query:"end_time" form:"end_time"`
}

// TallyList 获取所有账单
func TallyList(c echo.Context) error {
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.TallyCode, global.UserIdentityErr)
	}
	//缓存中获取
	val := global.Global.Redis.Get(global.Global.Ctx, global.TallyListKey+id).Val()
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
			return common.Fail(c, global.TallyCode, global.UserIdentityErr)
		}
		global.Global.Pool.Submit(func() {
			marshal, err := json.Marshal(list)
			if err != nil {
				return
			}
			val, err := global.Global.Redis.Set(global.Global.Ctx, global.TallyListKey+id, marshal, 0).Result()
			fmt.Println(val, err)
		})
		return common.Ok(c, list)
	}

}

// AddTallyLog 添加账单
func AddTallyLog(c echo.Context) error {
	t := new(List)
	userIdentity := utils.GetIdentity(c, "identity")
	if userIdentity == "" {
		return common.Fail(c, global.TallyCode, global.UserIdentityErr)
	}
	err := c.Bind(t)
	if err != nil {
		return common.Fail(c, global.TallyCode, global.ParseErr)
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
	global.Global.Pool.Submit(func() {
		val, err := global.Global.Redis.Del(global.Global.Ctx, userIdentity+"list").Result()
		global.Global.Redis.Del(global.Global.Ctx, userIdentity+strconv.Itoa(t.Category))
		fmt.Println("删除", val, err)
	})
	return common.Ok(c, nil)
}

// AllotKind 根据分类获取列表数据
func AllotKind(c echo.Context) error {
	category := c.QueryParam("category")
	userIdentity := utils.GetIdentity(c, "identity")
	if userIdentity == "" {
		return common.Fail(c, global.TallyCode, global.UserIdentityErr)
	}
	//解决数据不一致问题
	val := global.Global.Redis.Get(global.Global.Ctx, userIdentity+category).Val()
	//fmt.Println("val", val)
	if val != "" {
		var list []models.Tally
		err := json.Unmarshal([]byte(val), &list)
		if err != nil {
			return err
		}
		return common.Ok(c, list)
	} else {
		n, err := strconv.Atoi(category)
		if err != nil {
			return err
		}
		list := dao.GetTallyKind(userIdentity, n)
		global.Global.Pool.Submit(func() {
			marshal, err := json.Marshal(list)
			if err != nil {
				return
			}
			//存入redis
			global.Global.Redis.Set(global.Global.Ctx, userIdentity+category, marshal, 0)
		})
		return common.Ok(c, list)
	}
}

// DateList 按日期查询
func DateList(c echo.Context) error {
	//传时间戳
	t := new(Time)
	userIdentity := utils.GetIdentity(c, "identity")
	if userIdentity == "" {
		return common.Fail(c, global.TallyCode, global.UserIdentityErr)
	}
	err := c.Bind(t)
	if err != nil {
		fmt.Println(err)
		return common.Fail(c, global.TallyCode, global.ParseErr)
	}
	if t.StarTime > t.EndTime {
		return common.Fail(c, global.TallyCode, "时间错误")
	}
	star := time.Unix(t.StarTime, 0)
	end := time.Unix(t.EndTime, 0)
	val := global.Global.Redis.Get(global.Global.Ctx, userIdentity+strconv.FormatInt(t.EndTime, 10)).Val()
	fmt.Println("val", val)
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
				global.Global.Redis.Set(global.Global.Ctx, userIdentity+strconv.FormatInt(t.EndTime, 10), "null", time.Duration(utils.GetRandom(10))*time.Minute)
			}()
			return common.Fail(c, global.TallyCode, "查询失败")
		}
		global.Global.Pool.Submit(func() {
			marshal, err := json.Marshal(list)
			if err != nil {
				return
			}
			global.Global.Redis.Set(global.Global.Ctx, userIdentity+strconv.FormatInt(t.EndTime, 10), marshal, time.Duration(utils.GetRandom(10))*time.Minute)
		})
		return common.Ok(c, list)
	}
}

// BindKind 绑定分类
func BindKind(c echo.Context) error {
	category, err := strconv.Atoi(c.QueryParam("category"))
	if err != nil {
		fmt.Println(err)
		return common.Fail(c, global.TallyCode, global.ParseErr)
	}
	id := c.QueryParam("identity ")
	useIdentity := utils.GetIdentity(c, "identity")
	if useIdentity == "" {
		return common.Fail(c, global.TallyCode, global.UserIdentityErr)
	}
	if ok := dao.GetByKind(category); !ok {
		return common.Fail(c, global.TallyCode, "分类不存在")
	}
	err = dao.UpdateByKind(id, category)
	if err != nil {
		return common.Fail(c, global.TallyCode, "绑定失败")
	}

	return common.Ok(c, nil)
}

// Analysis 分析
func Analysis(c echo.Context) error {
	//获取
	list := make([]models.Tally, 0)
	err := c.Bind(&list)
	if err != nil {
		fmt.Println(err)
		return common.Fail(c, global.TallyCode, global.ParseErr)
	}
	fmt.Println(list)
	//使用星火大模型分析
	marshal, err := json.Marshal(list)
	if err != nil {
		return common.Fail(c, global.TallyCode, global.ParseErr)
	}
	//流式响应
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=stream.txt")

	conn := utils.GetConnect(marshal)

	var answer = ""
	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message error:", err)
			break
		}

		var data map[string]interface{}
		err1 := json.Unmarshal(msg, &data)
		if err1 != nil {
			fmt.Println("Error parsing JSON:", err)
			return common.Fail(c, global.TallyCode, global.ParseErr)
		}
		fmt.Println(string(msg))
		//解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			fmt.Println(data["payload"])
			return common.Fail(c, global.TallyCode, global.ParseErr)

		}
		status := choices["status"].(float64)
		fmt.Println(status)
		text := choices["text"].([]interface{})

		content := text[0].(map[string]interface{})["content"].(string)
		if status != 2 {
			//answer += content
			_, err = io.Copy(c.Response(), strings.NewReader(content))
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			fmt.Println("收到最终结果")
			answer += content
			usage := payload["usage"].(map[string]interface{})
			temp := usage["text"].(map[string]interface{})
			totalTokens := temp["total_tokens"].(float64)
			fmt.Println("total_tokens:", totalTokens)
			conn.Close()
			break
		}

	}
	//输出返回结果
	fmt.Println(answer)
	return nil
}
