package views

import (
	"Gin_Student_Management/models"
	"fmt"
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
	currentUserId := session.Get("userId").(uint)
	currentUserName := session.Get("userName")
	fmt.Println(currentUserName, currentUserId)
	isSuerAdmin := currentUserId == 1
	query := models.DB.Preload("Mentor")
	if !isSuerAdmin {
		query = query.Where("mentor_id = ?", currentUserId)
	}
	query.Find(&students)

	// 数据返回给前端页面
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Students": students,
		"Teachers": teachers,
		"User": gin.H{
			"ID":      currentUserId,
			"Name":    currentUserName,
			"IsSuper": isSuerAdmin,
		},
	})
}

// AddStudent 新增学生请求
func AddStudent(c *gin.Context) {
	session := sessions.Default(c)
	currentUserID := session.Get("userId").(uint)

	// 从前端获取数据
	name := c.PostForm("name")
	age := c.PostForm("age")
	gender := c.PostForm("gender")
	class := c.PostForm("class")

	formMentorID := c.PostForm("mentor_id")
	var finalMentorID uint

	ageInt, _ := strconv.Atoi(age)
	if currentUserID == 1 {
		mid, _ := strconv.ParseUint(formMentorID, 10, 64)
		finalMentorID = uint(mid)
	} else {
		finalMentorID = currentUserID
	}

	// 新建学生类组装信息
	newStudent := models.Student{
		Name:     name,
		Age:      ageInt,
		Gender:   gender,
		Class:    class,
		MentorID: finalMentorID,
		JoinDate: time.Now(),
	}

	err := models.DB.Create(&newStudent).Error
	if err != nil {
		c.String(http.StatusBadRequest, "添加失败")
		return
	}
	// 添加成功，重定向回主页
	c.Redirect(http.StatusFound, "/")
}

func DeleteStudent(c *gin.Context) {
	session := sessions.Default(c)
	currentUserID := session.Get("userId").(uint)

	id := c.Param("id")
	var student models.Student

	err := models.DB.Where("id = ?", id).First(&student).Error
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	if currentUserID != 1 && student.MentorID != currentUserID {
		c.String(http.StatusForbidden, "无权删除")
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

func EditStudentPage(c *gin.Context) {
	id := c.Param("id")
	session := sessions.Default(c)
	currentUserID := session.Get("userId").(uint)

	var student models.Student
	var teachers []models.Teacher

	if err := models.DB.Preload("Mentor").First(&student, id).Error; err != nil {
		c.String(http.StatusNotFound, "学生不存在")
		return
	}

	if currentUserID != 1 && student.MentorID != currentUserID {
		c.String(http.StatusForbidden, "你无权编辑该学生")
		return
	}

	models.DB.Find(&teachers)

	c.HTML(http.StatusOK, "edit.tmpl", gin.H{
		"Student":  student,
		"Teachers": teachers,
		"User": gin.H{
			"IsSuper": currentUserID == 1,
		},
	})
}

func UpdateStudent(c *gin.Context) {
	session := sessions.Default(c)
	currentUserID := session.Get("userId").(uint)

	id := c.PostForm("id")

	var student models.Student
	if err := models.DB.First(&student, id).Error; err != nil {
		c.String(http.StatusNotFound, "Error")
		return
	}

	if currentUserID != 1 && student.MentorID != currentUserID {
		c.String(http.StatusForbidden, "你无权修改该学生")
		return
	}

	name := c.PostForm("name")
	age, _ := strconv.Atoi(c.PostForm("age"))
	gender := c.PostForm("gender")
	class := c.PostForm("class")

	var newMentorID uint
	if currentUserID == 1 {
		mid, _ := strconv.ParseUint(c.PostForm("mentor_id"), 10, 64)
		newMentorID = uint(mid)
	} else {
		newMentorID = student.MentorID
	}

	updateData := map[string]interface{}{
		"name":      name,
		"age":       age,
		"gender":    gender,
		"class":     class,
		"mentor_id": newMentorID,
	}

	if err := models.DB.Model(&student).Where("id = ?", id).Update(updateData).Error; err != nil {
		c.String(http.StatusInternalServerError, "更新失败")
		return
	}

	c.Redirect(http.StatusFound, "/")
}
