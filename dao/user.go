package dao

import (
	"Tally/global"
	"Tally/models"
	"fmt"
)

func GetUserById(username, password string) *models.User {
	//
	user := new(models.User)
	err := global.Global.Mysql.Where("username=? and password=? ", username, password).Take(user).Error
	if err != nil {
		return nil
	}
	return user
}

func GetPhone(phone string) bool {
	user := new(models.User)
	err := global.Global.Mysql.Where("phone=?", phone).Take(user).Error
	if err != nil {
		return false
	}
	return true
}

func InsertUser(user *models.User) error {
	err := global.Global.Mysql.Create(user).Error
	return err
}

func GetUserByUsername(username string) bool {
	user := new(models.User)
	err := global.Global.Mysql.Where("username=?", username).Take(user).Error
	if err != nil {
		fmt.Println("sql查询错误" + err.Error())
		return false
	}
	return true
}
