package dao

import (
	"Tally/global"
	"Tally/models"
)

func InsertFollow(follow *models.Follow) error {
	return global.Global.Mysql.Create(follow).Error
}
