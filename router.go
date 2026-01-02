package main

import (
	"Gin_Student_Management/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func initRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"ok": "ok",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		students := []models.Student{
			{ID: 1, Name: "张三", Age: 20, Gender: "男", Class: "计算机一班", JoinDate: "2023-09-01"},
			{ID: 2, Name: "李四", Age: 19, Gender: "女", Class: "计算机二班", JoinDate: "2023-09-01"},
		}
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"Students": students,
		})
	})
}
