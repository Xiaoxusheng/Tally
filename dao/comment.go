package dao

import (
	"Tally/global"
	"Tally/models"
)

func GetCommentList(blogId string) ([]*models.Comment, error) {
	comment := make([]*models.Comment, 0)
	return comment, global.Global.Mysql.Where("blog_id=?", blogId).Find(comment).Error
}
