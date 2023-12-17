package config

import (
	"Tally/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMysql() {
	dsn := Config.Mysql.Username + ":" + Config.Mysql.Password + "@tcp(" + Config.Mysql.Url + ")/" + Config.Mysql.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	db.Debug()
	global.Global.Log.Info("mysql初始化成功")
	global.Global.Mysql = db
	//建表
	//err = global.Global.Mysql.AutoMigrate(&models.Blog{})
	//err = global.Global.Mysql.AutoMigrate(&models.Collect{})
	//if err != nil {
	//	global.Global.Log.Info(err)
	//	return
	//}
	//global.Global.Mysql.AutoMigrate(&models.Comment{})
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
