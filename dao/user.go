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

func GetByPwdIdentity(id, pwd string) *models.User {
	user := new(models.User)
	err := global.Global.Mysql.Where("identity=? and password=? ", id, pwd).Take(user).Error
	if err != nil {
		return nil
	}
	return user
}

func UpdatePwd(id, pwd string) error {
	user := new(models.User)
	err := global.Global.Mysql.Model(user).Where("identity=?", id).Update("password", pwd).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserByIdentity(id string) *models.User {
	//
	user := new(models.User)
	err := global.Global.Mysql.Where("identity=? ", id).Take(user).Error
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
		global.Global.Log.Warn("sql查询错误" + err.Error())
		return false
	}
	return true
}

func DeleteUser(id string) error {
	user := new(models.User)
	return global.Global.Mysql.Where("identity=?", id).Delete(user).Error
}

// UpdateAll 全量更新
func UpdateAll(id string, user *global.User) error {
	fmt.Println(user)
	users := new(models.User)
	var str, s string
	if user.GithubId != "" {
		str = "GithubId"
		users.GithubId = user.GithubId
	}
	if user.Phone != "" {
		s = "Phone"
		users.Phone = user.Phone
	}
	return global.Global.Mysql.Model(users).Select(str, s).Where("identity=?", id).Updates(models.User{
		Phone:    user.Phone,
		GithubId: user.GithubId,
	}).Error
}
