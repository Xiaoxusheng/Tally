package dao

import (
	"Tally/global"
	"Tally/models"
)

func GetKindList() []*models.Kind {
	kind := make([]*models.Kind, 0)
	err := global.Global.Mysql.Find(&kind).Error
	if err != nil {
		return nil
	}
	return kind
}

func GetByKind(s int) bool {
	kind := new(models.Kind)
	err := global.Global.Mysql.Where("serial_number=?", s).Take(kind).Error
	if err != nil {
		return false
	}
	return true
}

func InsertKind(u *models.Kind) error {
	err := global.Global.Mysql.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}
