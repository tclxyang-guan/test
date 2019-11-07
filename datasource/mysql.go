package datasource

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"os"
	"test/config"
	"test/models"
	"time"
)

var mysqldb *gorm.DB

func init() {
	var sql = config.Sysconfig.Mysql.Name + ":" + config.Sysconfig.Mysql.Password + "@(" + config.Sysconfig.Mysql.Host + ":" + config.Sysconfig.Mysql.Port + ")/" + config.Sysconfig.Mysql.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", sql)
	if err != nil {
		log.Fatal("连接数据库失败")
	}
	db.DB().SetMaxIdleConns(10)
	db.AutoMigrate(
		models.Test{},
	)
	db.LogMode(true)
	logger := log.New()
	if config.Sysconfig.Basic.IsLocal {
		logger.Out = os.Stdout
	} else {
		logger.Out = NewLogFile()
		db.SetLogger(logger)
	}

	mysqldb = db
}
func GetDB() *gorm.DB {
	return mysqldb
}

//请求日志以及sql日志
func NewLogFile() *os.File {
	filename := config.Sysconfig.Basic.ReqLogPath + "req_" + time.Now().Format("20060102") + ".txt"
	//打开一个输出文件，如果重新启动服务器，它将追加到今天的文件中
	f, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	return f
}
