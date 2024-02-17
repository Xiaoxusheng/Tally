package dao

import (
	"Tally/global"
	"Tally/models"
)

func InsertBlog(blog *models.Blog) error {
	err := global.Global.Mysql.Create(blog).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateLikes(id string, n string) error {
	blog := new(models.Blog)
	err := global.Global.Mysql.Model(blog).Where("identity=?", id).Update("likes", n).Error
	if err != nil {
		return err
	}
	return nil
}

func GetIdByBlog(id string) string {
	blog := new(models.Blog)
	err := global.Global.Mysql.Where("Identity=?", id).Take(blog).Error
	if err != nil {
		return ""
	}
	return blog.UserIdentity
}

func DeleteBlogByIdentity(id string) error {
	blog := new(models.Blog)
	return global.Global.Mysql.Where("identity=?", id).Delete(blog).Error
}

func UpdateStatus(id string, status int) error {
	blog := new(models.Blog)
	return global.Global.Mysql.Model(blog).Where("identity=?", id).Update("is_hide", status).Error
}

func GetBlogListById(id string) ([]*models.Blog, error) {
	blog := make([]*models.Blog, 0)
	err := global.Global.Mysql.Where("identity=?", id).Find(&blog).Error
	return blog, err
}
