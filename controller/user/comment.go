package user

import (
	"Tally/common"
	"Tally/global"
	"Tally/utils"
	"github.com/labstack/echo/v4"
)

/*
评论模块
*/

// PushComment 发表评论
func PushComment(c echo.Context) error {
	//获取

	//异步检验内容

	//发布

	return common.Ok(c, nil)
}

// DeleteComment 删除评论
func DeleteComment(c echo.Context) error {
	return common.Ok(c, nil)

}

// GetCommentList 获取评论
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
