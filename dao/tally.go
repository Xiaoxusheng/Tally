package dao

import (
	"Tally/global"
	"Tally/models"
	"fmt"
)

func GetTallyList(id string) []*models.Tally {
	list := make([]*models.Tally, 0)
	err := global.Global.Mysql.Where("user_identity=?", id).Find(&list).Error
	if err != nil {
		panic(err)
		return nil
	}
	fmt.Println("list", list)
	return list
}

func InsertTally(u *models.Tally) error {
	err := global.Global.Mysql.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}

func GetTallyKind(identity string, n int) []*models.Tally {
	list := make([]*models.Tally, 0)
	err := global.Global.Mysql.Where("category=? and user_identity=?", n, identity).Find(&list).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return list
}

func GetByTime(star, end string) []*models.Tally {
	list := make([]*models.Tally, 0)
	err := global.Global.Mysql.Where("created_at>=? and created_at<=?", star, end).Find(&list).Error
	if err != nil {
		return nil
	}
	return list
}

func UpdateByKind(id string, kind int) error {
	list := new(models.Tally)
	err := global.Global.Mysql.Model(list).Where("identity=?", id).Update("category", kind).Error
	if err != nil {
		return err
	}
	return nil
}

func GetLikeList(s string) []*models.Tally {
	list := make([]*models.Tally, 0)
	err := global.Global.Mysql.Where("remark like ?", "%"+s+"%").Find(&list).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return list
}

func GetListById(id string) []*models.Tally {
	list := make([]*models.Tally, 0)
	err := global.Global.Mysql.Where("identity like ?", "%"+id+"%").Find(&list).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return list
}
