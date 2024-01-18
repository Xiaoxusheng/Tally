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
评论模块
*/

// PushComment 发表评论
func PushComment(c echo.Context) error {
	//获取

	//todo 送进消息队列中,进行校验

	return common.Ok(c, nil)
}

// DeleteComment 删除评论
func DeleteComment(c echo.Context) error {

	return common.Ok(c, nil)

}

// GetCommentList 获取评论列表
func GetCommentList(c echo.Context) error {
	blogId := c.QueryParam("blog_id")
	if blogId == "" {
		return common.Fail(c, global.CommentCode, global.QueryErr)
	}
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.TallyCode, global.UserIdentityErr)
	}
	//获取缓存

	return common.Ok(c, nil)

}

// CommentList 评论列表
func CommentList(c echo.Context) error {
	blogId := c.QueryParam("blog_id")
	if blogId == "" {
		return common.Fail(c, global.CommentCode, global.QueryErr)
	}
	//判断blog是否存在
	value := global.Global.Redis.SIsMember(global.Global.Ctx, global.BlogText+":IdList", blogId).Val()
	if !value {
		return common.Fail(c, global.LikesCode, global.BlogNotFound)
	}

	val := global.Global.Redis.SMembers(global.Global.Ctx, global.CommentList+blogId).Val()
	list := make([]*models.Comment, 0, global.Global.Redis.SCard(global.Global.Ctx, global.CommentList+blogId).Val())
	if len(val) != global.Fail {
		for _, res := range val {
			comment := new(models.Comment)
			err := json.Unmarshal([]byte(res), comment)
			if err != nil {
				return common.Fail(c, global.CommentCode, global.MarshalErr)
			}
			list = append(list, comment)
		}
		return common.Ok(c, list)
	} else {
		lists, err := dao.GetCommentList(blogId)
		if err != nil {
			return common.Fail(c, global.CommentCode, global.GetCommentListErr)
		}
		//写进缓存
		global.Global.Pool.Submit(func() {
			for i := 0; i < len(lists); i++ {
				marshal, err := json.Marshal(lists[i])
				if err != nil {
					global.Global.Log.Error("marshal err:", err)
				}
				global.Global.Redis.SAdd(global.Global.Ctx, global.CommentList+blogId, marshal).Val()
			}
		})

		return common.Ok(c, lists)
	}

}
