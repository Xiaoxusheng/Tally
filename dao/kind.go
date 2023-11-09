package dao

import (
	"Tally/global"
	"Tally/models"
)

func GetKind() []*models.Kind {
	kind := make([]*models.Kind, 0)
	err := global.Global.Mysql.Find(kind).Error
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
