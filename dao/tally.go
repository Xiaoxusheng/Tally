package dao

import (
	"Tally/global"
	"Tally/models"
)

func GetTallyList(id string) *models.Tally {
	list := new(models.Tally)
	err := global.Global.Mysql.Where("user_identity=?", id).Find(list).Error
	if err != nil {
		return nil
	}
	return list
}

func InsertTally(u *models.Tally) error {
	err := global.Global.Mysql.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}
