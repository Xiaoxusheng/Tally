package dao

import (
	"Tally/global"
	"Tally/models"
)

func InsertBlogCollect(collect *models.Collect) error {
	return global.Global.Mysql.Create(collect).Error
}

// DeleteBlogCollect 删除
func DeleteBlogCollect(blogId string) error {
	collect := new(models.Collect)
	return global.Global.Mysql.Where("blog_id=?", blogId).Delete(collect).Error
}

func UpdateBlogCollect(blogId string) error {
	collect := new(models.Collect)
	global.Global.Log.Info("修改")
	return global.Global.Mysql.Model(collect).Unscoped().Where("blog_id=?", blogId).Update("deleted_at", nil).Error
}
