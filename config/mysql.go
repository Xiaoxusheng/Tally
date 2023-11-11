package config

import (
	"Tally/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func InitMysql() {
	dsn := Config.Mysql.Username + ":" + Config.Mysql.Password + "@tcp(127.0.0.1:3306)/" + Config.Mysql.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	db.Debug()
	log.Println("mysql初始化成功")
	global.Global.Mysql = db
	//f := sync.Once{}
	//f.Do(
	//	func() {
	//		err = global.Global.Mysql.AutoMigrate(&models.User{})
	//		err = global.Global.Mysql.AutoMigrate(&models.Tally{})
	//		err = global.Global.Mysql.AutoMigrate(&models.Kind{})
	//	})
	//if err != nil {
	//	panic(err)
	//}
}
