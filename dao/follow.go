package dao

import (
	"Tally/global"
	"Tally/models"
)

func InsertFollow(follow *models.Follow) error {
	return global.Global.Mysql.Create(follow).Error
}

func GetFollowList(id string) []*models.Follow {
	follow := make([]*models.Follow, 0)
	err := global.Global.Mysql.Where("user_id=?", id).Find(&follow).Error
	if err != nil {
		global.Global.Log.Warn(err)
	}
	return follow

}
