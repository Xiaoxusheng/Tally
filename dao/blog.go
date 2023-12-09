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

func UpDateLikes(id string, n string) error {
	blog := new(models.Blog)
	err := global.Global.Mysql.Model(blog).Where("identity=?", id).Update("likes", n).Error
	if err != nil {
		return err
	}
	return nil
}
