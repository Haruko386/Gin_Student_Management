package main

import (
	"Gin_Student_Management/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func main() {
	// init
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.Static("assets", "./assets")

	err := models.InitMysql()
	if err != nil {
		return
	}

	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(models.DB)
	models.DB.AutoMigrate(&models.Student{}, &models.PaperList{}, &models.Teacher{})

	initRoutes(r)

	err = r.Run(":8081")
	if err != nil {
		return
	}
}
