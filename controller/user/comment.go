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
	comment := new(global.Comment)
	err := c.Bind(comment)
	if err != nil {
		return common.Fail(c, global.CommentCode, global.ParseErr)
	}
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.LikesCode, global.UserIdentityErr)
	}
	//判断博客id是否存在
	//判断博客id是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.BlogText+":IdList", comment.BlogID).Val()
	if !val {
		return common.Fail(c, global.LikesCode, global.BlogNotFound)
	}
	com := utils.GetUidV4()
	//判断父评论是否存在
	// 存在父评论
	if comment.ParentID != "" {
		//判断是否存在
		if !global.Global.Redis.SIsMember(global.Global.Ctx, global.CommentId+com, comment.ParentID).Val() {
			return common.Fail(c, global.CommentCode, global.CommentParentIdNotFound)
		}
	}
	//不存在
	//将生成的的评论唯一id添加进set中
	uid := utils.GetUidV4()
	res := global.Global.Redis.SAdd(global.Global.Ctx, global.CommentId+com, uid).Val()
	if res == global.Fail {
		return common.Fail(c, global.CommentCode, global.CommentFail)
	}
	//todo 送进消息队列中,进行校验

	//添加进数据库
	comments := &models.Comment{
		Identity:     com,
		UserIdentity: id,
		BlogID:       comment.BlogID,
		ParentID:     comment.ParentID,
		Text:         comment.Text,
		Ip:           c.RealIP(),
		IsTop:        false,
		ViolateRule:  false,
	}
	err = dao.InsertComment(comments)
	if err != nil {
		return common.Fail(c, global.CommentCode, global.CommentFail)
	}
	//写进缓存
	global.Global.Pool.Submit(func() {
		marshal, err := json.Marshal(comments)
		if err != nil {
			global.Global.Log.Error("marshal err:", err)
		}
		//添加进缓存
		global.Global.Redis.SAdd(global.Global.Ctx, global.CommentList+comment.BlogID, marshal).Val()
	})

	return common.Ok(c, nil)
}

// DeleteComment 删除评论
func DeleteComment(c echo.Context) error {
	//
	delComment := new(global.DelComment)
	err := c.Bind(delComment)
	if err != nil {
		return common.Fail(c, global.CommentCode, global.ParseErr)
	}
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.CommentCode, global.UserIdentityErr)
	}
	//判断身份

	//删除
	err = dao.DelCommentById(delComment.CommentId)
	if err != nil {
		return common.Fail(c, global.CommentCode, global.DelCommentFail)
	}
	global.Global.Pool.Submit(func() {
		comment, err := dao.GetCommentById(delComment.CommentId)
		if err != nil {
			return
		}
		marshal, err := json.Marshal(comment)
		if err != nil {
			return
		}
		//集合
		global.Global.Redis.SRem(global.Global.Ctx, global.CommentList+delComment.BlogID, marshal).Val()
	})
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
				//集合
				global.Global.Redis.SAdd(global.Global.Ctx, global.CommentList+blogId, marshal).Val()
			}
		})

		return common.Ok(c, lists)
	}

}
