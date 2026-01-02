package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitMysql() (err error) {
	dsn := "root:password@(127.0.0.1:3306)/student_management_gorm?charset=utf8&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = DB.DB().Ping()
	return err
}

type Student struct {
	ID        uint
	StudentID uint
	Name      string
	Age       int
	Gender    string
	Class     string
	JoinDate  string
}

type Teacher struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
}

type PaperList struct {
	Author  Student `json:"author"`
	Journal string  `json:"journal"`
	Title   string  `json:"title"`
}
