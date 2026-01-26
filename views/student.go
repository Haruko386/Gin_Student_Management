package views

import (
	"Gin_Student_Management/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetStudentList 获取学生列表
func GetStudentList(c *gin.Context) {
	var students []models.Student
	var teachers []models.Teacher

	models.DB.Preload("Mentor").Find(&students)
	models.DB.Find(&teachers)

	session := sessions.Default(c)
	currentUserName := session.Get("userName")
	isAdmin := session.Get("isAdmin")

	// 数据返回给前端页面
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Students": students,
		"Teachers": teachers,
		"User": gin.H{
			"Name":    currentUserName,
			"IsAdmin": isAdmin,
		},
	})
}

// AddStudent 新增学生请求
func AddStudent(c *gin.Context) {
	// 从前端获取数据
	name := c.PostForm("name")
	age := c.PostForm("age")
	gender := c.PostForm("gender")
	class := c.PostForm("class")
	mentorID := c.PostForm("mentor_id")

	ageInt, _ := strconv.Atoi(age)
	mentorIDUint, _ := strconv.ParseUint(mentorID, 10, 64)

	// 新建学生类组装信息
	newStudent := models.Student{
		Name:     name,
		Age:      ageInt,
		Gender:   gender,
		Class:    class,
		MentorID: uint(mentorIDUint),
		JoinDate: time.Now(),
	}

	err := models.DB.Create(&newStudent).Error
	if err != nil {
		c.String(http.StatusBadRequest, "添加失败")
		return
	}
	println("dwadwa")
	// 添加成功，重定向回主页
	c.Redirect(http.StatusFound, "/")
}

func DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	var student models.Student

	err := models.DB.Where("id = ?", id).First(&student).Error
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	// 删除
	err = models.DB.Delete(&student).Error
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.Redirect(http.StatusFound, "/")
}

func UpdateStudent(c *gin.Context) {

}
