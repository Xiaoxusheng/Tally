package config

import (
	"Tally/global"
	"Tally/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func InitMysql() {
	dsn := Config.Mysql.Username + ":" + Config.Mysql.Password + "@tcp(127.0.0.1:3306)/" + Config.Mysql.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Debug()
	log.Println("mysql初始化成功")
	global.Global.Mysql = db
	global.Global.Mysql.AutoMigrate(&models.User{})
}
