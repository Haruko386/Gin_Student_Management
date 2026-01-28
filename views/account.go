package views

import (
	"Gin_Student_Management/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	id := c.PostForm("username")
	password := c.PostForm("password")
	role := c.PostForm("role")

	session := sessions.Default(c)

	if role == "teacher" {
		var teacher models.Teacher
		if err := models.DB.Where("id = ? AND password = ?", id, password).First(&teacher).Error; err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{"Error": "Login Failed: Invalid ID or Password"})
			return
		}
		// 老师登录成功
		session.Set("userId", teacher.ID)
		session.Set("userName", teacher.Name)
		session.Set("role", "teacher")          // 记录角色
		session.Set("isAdmin", teacher.IsAdmin) // 只有老师有这属性
		session.Save()
		if teacher.ID == 1 {
			c.Redirect(http.StatusFound, "/admin")
		} else {
			c.Redirect(http.StatusFound, "/")
		}

	} else {
		// 学生登录逻辑
		var student models.Student
		if err := models.DB.Where("id = ? AND password = ?", id, password).First(&student).Error; err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{"Error": "Student Login Failed"})
			return
		}
		// 学生登录成功
		session.Set("userId", student.ID)
		session.Set("userName", student.Name)
		session.Set("role", "student")
		session.Set("isAdmin", false)
		session.Save()
		c.Redirect(http.StatusFound, "/student/dashboard")
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		return
	}
	c.Redirect(http.StatusFound, "/login")
}

func UpdatePassword(c *gin.Context) {
	oldPwd := c.PostForm("old_password")
	newPwd := c.PostForm("new_password")
	confirmPwd := c.PostForm("confirm_password")

	if newPwd != confirmPwd {
		c.String(http.StatusBadRequest, "两次输入密码不一致")
		return
	}

	session := sessions.Default(c)
	userId := session.Get("userId")
	if userId == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var teacher models.Teacher
	if err := models.DB.Where("id = ? AND password = ?", userId, oldPwd).First(&teacher).Error; err != nil {
		c.String(http.StatusBadRequest, "旧密码错误")
		return
	}

	models.DB.Model(&teacher).Update("password", newPwd)
	c.String(http.StatusOK, "密码修改成功")
}

func LoginPage(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("userId") != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})
}
