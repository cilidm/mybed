package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"mybedv2/app/helper/util/pathdir"
	"mybedv2/conf"
	"os"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func Instance() *gorm.DB {
	once.Do(func() {
		var err error
		if conf.Setting.DBType == "mysql" {
			db, err = gorm.Open(conf.Setting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				conf.Setting.DBUser,
				conf.Setting.DBPwd,
				conf.Setting.DBHost,
				conf.Setting.DBTableName))
		} else if conf.Setting.DBType == "sqlite3" {
			db, err = gorm.Open(conf.Setting.DBType, conf.Setting.DBPath)
			if err != nil {
				if err := pathdir.CreateFile(conf.Setting.DBPath); err != nil {
					log.Println(err)
				}
			}
		}
		if err != nil {
			log.Println("数据库连接失败:【" + err.Error() + "】")
			os.Exit(-1)
		}

		db.SingularTable(true)
		if conf.Setting.RunMode != "release" {
			db.LogMode(true)
		}
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
		initModel() //校验数据表
		initAdmin() //校验admin及必填数据
	})
	return db
}
