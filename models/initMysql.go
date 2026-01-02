package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitMysql() (err error) {
	dsn := "root:password@(127.0.0.1:3306)/gorm_demo?charset=utf8&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = DB.DB().Ping()
	return err
}

type Student struct {
	Id    int
	Name  string
	Grade string
}

type PaperList struct {
	Author  Student
	Journal string
	Title   string
}
