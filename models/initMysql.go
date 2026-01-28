package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var DB *gorm.DB

func InitMysql() (err error) {
	dsn := "root:password@(127.0.0.1:3306)/student_management_gorm?charset=utf8&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	DB.AutoMigrate(&Student{}, &PaperList{}, &Teacher{}, &Paper{})

	var count int
	DB.Model(&Teacher{}).Where("id = ?", 1).Count(&count)
	if count == 0 {
		superAdmin := Teacher{
			Model:    gorm.Model{ID: 1},
			Name:     "SuperAdmin",
			Password: "password",
			IsAdmin:  true,
		}
		DB.Create(&superAdmin)
	}

	err = DB.DB().Ping()
	return err
}

type Student struct {
	gorm.Model
	Password  string    `json:"password"`
	StudentID uint      `json:"student_id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	Class     string    `json:"class"`
	JoinDate  time.Time `json:"join_date"`

	MentorID uint    `json:"mentor"`
	Mentor   Teacher `gorm:"foreignKey:MentorID"`
}

type Teacher struct {
	gorm.Model
	Password string `json:"password"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"isAdmin"`
}

type Paper struct {
	gorm.Model
	StudentID uint   `json:"studentId"` // 外键：属于哪个学生
	Title     string `json:"title"`
	CoverPath string `json:"coverPath"` // 封面图片路径
	FilePath  string `json:"filePath"`  // 论文文件路径
	Status    string `json:"status"`    // 状态：已上传、待审核等
}

type PaperList struct {
	gorm.Model
	AuthorID uint   `json:"author_id"`
	Journal  string `json:"journal"`
	Title    string `json:"title"`
	Storage  string `json:"storage"`
}
