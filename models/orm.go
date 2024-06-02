package models

import (
	"errors"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

var orm *gorm.DB
var GDRSDns string //主数据库
func NewOrm() (*gorm.DB, error) {

	if orm == nil {
		var err error
		if orm, err = gorm.Open(mysql.Open(GDRSDns), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
		}); err != nil {
			log.Println("conn mysql error!", err)
			return nil, errors.New("读取数据库错误")
		}
		sqlDB, _ := orm.DB()
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	return orm, nil
}
