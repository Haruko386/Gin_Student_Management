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

	var teacher models.Teacher

	if err := models.DB.Where("id = ? AND password = ?", id, password).First(&teacher).Error; err != nil {
		c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{"error": "登陆失败，ID或密码错误"})
		return
	}

	session := sessions.Default(c)
	session.Set("userId", teacher.ID)
	session.Set("username", teacher.Name)
	session.Set("isAdmin", teacher.IsAdmin)
	session.Save()

	if teacher.IsAdmin {
		c.Redirect(http.StatusFound, "/admin")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

func UpdatePassword(c *gin.Context) {
	oldPwd := c.PostForm("oldPwd")
	newPwd := c.PostForm("newPwd")
	confirmPwd := c.PostForm("confirm")

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
