package config

import (
	"Tally/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

func InitMysql() {
	once := sync.Once{}
	once.Do(
		func() {
			newLogger := logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold:             time.Second,   // Slow SQL threshold
					LogLevel:                  logger.Silent, // Log level
					IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
					ParameterizedQueries:      true,          // Don't include params in the SQL log
					Colorful:                  false,         // Disable color
				},
			)
			dsn := Config.Mysql.Username + ":" + Config.Mysql.Password + "@tcp(" + Config.Mysql.Url + ")/" + Config.Mysql.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger:      newLogger,
				PrepareStmt: true,
			})
			if err != nil {
				panic(err)
			}
			db.Debug()
			global.Global.Log.Info("mysql初始化成功")
			global.Global.Mysql = db
		})
	//建表
	//err = global.Global.Mysql.AutoMigrate(&models.Follow{})
	//err = global.Global.Mysql.AutoMigrate(&models.Blog{})
	//err = global.Global.Mysql.AutoMigrate(&models.Comment{})
	//err = global.Global.Mysql.AutoMigrate(&models.Collect{})
	//err = global.Global.Mysql.AutoMigrate(&models.Collect{})
	//if err != nil {
	//	global.Global.Log.Info(err)
	//	return
	//}
	//err = global.Global.Mysql.AutoMigrate(&models.Comment{})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
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
