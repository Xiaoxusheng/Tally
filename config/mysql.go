package config

import (
	"Tally/global"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var (
	onc sync.Once
	db  *gorm.DB
	err error
)

func InitMysql() {
	onc.Do(
		func() {
			newLogger := logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold:             time.Second, // Slow SQL threshold
					LogLevel:                  logger.Info, // Log level
					IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
					ParameterizedQueries:      true,        // Don't include params in the SQL log
					Colorful:                  false,       // Disable color
				},
			)
			dsn := Config.Mysql.Username + ":" + Config.Mysql.Password + "@tcp(" + Config.Mysql.Url + ")/" + Config.Mysql.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger:      newLogger,
				PrepareStmt: true,
			})
			if err != nil {
				panic(err)
			}
			db.Debug()
			global.Global.Log.Info("mysql连接成功！")
			mysqlDB, _ := db.DB()

			// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
			mysqlDB.SetMaxIdleConns(Config.Mysql.MaxIdleCons)

			// SetMaxOpenConns 设置打开数据库连接的最大数量。
			mysqlDB.SetMaxOpenConns(Config.Mysql.MaxOpenCons)

			// SetConnMaxLifetime 设置了连接可复用的最大时间。
			mysqlDB.SetConnMaxLifetime(time.Minute * time.Duration(Config.Mysql.ConnMaxLifetime))
			global.Global.Mysql = db
		})
	//建表
	//err := global.Global.Mysql.AutoMigrate(&models.Follow{})
	//err = global.Global.Mysql.AutoMigrate(&models.Blog{})
	//err = global.Global.Mysql.AutoMigrate(&models.User{})
	//err = global.Global.Mysql.AutoMigrate(&models.Kind{})
	//err = global.Global.Mysql.AutoMigrate(&models.Tally{})
	//err = global.Global.Mysql.AutoMigrate(&models.Comment{})
	//err = global.Global.Mysql.AutoMigrate(&models.Collect{})
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

func InitCasBin() {
	sql, err := gormadapter.NewAdapter("mysql", Config.Mysql.Username+":"+Config.Mysql.Password+"@tcp("+Config.Mysql.Url+")/"+Config.Mysql.Database+"?charset=utf8mb4&parseTime=True&loc=Local", true) // Your driver and data source.
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer("./global/model.conf", sql)
	if err != nil {
		panic(err)
	}
	// Or you can use an existing DB "abc" like this:
	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	// a := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/abc", true)

	// Load the policy from DB.
	err = e.LoadPolicy()
	if err != nil {
		panic(e)
	}

	global.Global.CasBin = e
	// Check the permission.

}
