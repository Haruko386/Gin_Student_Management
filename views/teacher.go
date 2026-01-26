package views

import (
	"Gin_Student_Management/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AdminDashBoard 渲染后台页面
func AdminDashBoard(c *gin.Context) {
	var teacher []models.Teacher
	models.DB.Find(&teacher)

	c.HTML(http.StatusOK, "admin.tmpl", gin.H{
		"Teachers": teacher,
	})
}

// AddTeacher 删除教师
func AddTeacher(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	isAdminVal := c.PostForm("isAdmin")

	isAdmin := false
	if isAdminVal == "true" {
		isAdmin = true
	}

	newTeacher := models.Teacher{
		Name:     name,
		Password: password,
		IsAdmin:  isAdmin,
	}

	if err := models.DB.Create(&newTeacher).Error; err != nil {
		c.String(http.StatusBadRequest, "创建失败")
		return
	}
	c.Redirect(http.StatusFound, "/admin")
}

func DeleteTeacher(c *gin.Context) {
	id := c.Param("id")

	if id == "1" {
		c.String(http.StatusForbidden, "禁止删除超级管理员")
		return
	}

	models.DB.Where("id = ?", id).Delete(&models.Teacher{})
	c.Redirect(http.StatusFound, "/admin")
}

func GrantAdmin(c *gin.Context) {
	id := c.Param("id")
	models.DB.Model(&models.Teacher{}).Where("id = ?", id).Update("is_admin", true)
	c.Redirect(http.StatusFound, "/admin")
}

func RevokeAdmin(c *gin.Context) {
	id := c.Param("id")
	if id == "1" {
		c.String(http.StatusForbidden, "禁止撤销超级管理员权限")
		return
	}

	models.DB.Model(&models.Teacher{}).Where("id = ?", id).Update("is_admin", false)
	c.Redirect(http.StatusFound, "/admin")
}
