package dao

import (
	"Tally/global"
	"Tally/models"
)

func GetCommentList(blogId string) ([]*models.Comment, error) {
	comment := make([]*models.Comment, 0)
	return comment, global.Global.Mysql.Where("blog_id=?", blogId).Find(comment).Error
}

func InsertComment(comment *models.Comment) error {
	return global.Global.Mysql.Create(comment).Error
}

func DelCommentById(id string) error {
	comment := new(models.Comment)
	return global.Global.Mysql.Where("identity=?", id).Delete(comment).Error
}

func GetCommentById(id string) (*models.Comment, error) {
	comment := new(models.Comment)
	return comment, global.Global.Mysql.Where("Identity=?", id).Take(comment).Error
}
