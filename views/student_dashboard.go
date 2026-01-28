package views

import (
	"Gin_Student_Management/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StudentDashboard(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userId")

	if userID == nil || session.Get("role") != "student" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 获取学生信息
	var student models.Student
	models.DB.Preload("Mentor").First(&student, userID)

	//获取学生论文
	var paper models.Paper
	models.DB.Where("student_id = ?", student.ID).First(&paper)

	c.HTML(http.StatusOK, "student_dashboard.tmpl", gin.H{
		"Student": student,
		"Paper":   paper,
	})
}

func UploadThesis(c *gin.Context) {
	session := sessions.Default(c)
	studentID := session.Get("userId")

	file, err := c.FormFile("thesis_file")
	if err != nil {
		c.String(http.StatusBadRequest, "upload file failed")
		return
	}

	cover, err2 := c.FormFile("cover_image")
	if err2 != nil {
		c.String(http.StatusBadRequest, "cover image file failed")
		return
	}

	// 上传处理
	dstFile := "./static/uploads/thesis_" + file.Filename
	dstCover := "./static/uploads/cover_" + cover.Filename

	err = c.SaveUploadedFile(file, dstFile)
	if err != nil {
		return
	}
	err = c.SaveUploadedFile(file, dstCover)
	if err != nil {
		return
	}

	// 更新数据库
	var paper models.Paper
	if models.DB.Where("student_id = ?", studentID).First(&paper).RecordNotFound() {
		paper = models.Paper{
			StudentID: studentID.(uint),
			Title:     c.PostForm("Title"),
			FilePath:  "/static/uploads/thesis_" + file.Filename,
			CoverPath: "/static/uploads/cover_" + cover.Filename,
			Status:    "Submitted",
		}
		models.DB.Create(&paper)
	} else {
		paper.Title = c.PostForm("title")
		paper.FilePath = "/static/uploads/thesis_" + file.Filename
		paper.CoverPath = "/static/uploads/cover_" + cover.Filename
		paper.Status = "Updated"
		models.DB.Save(&paper)
	}

	c.Redirect(http.StatusFound, "/student/dashboard")
}
