package main

import (
	"Gin_Student_Management/views"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func initRoutes(r *gin.Engine) {
	// 初始化session中间件
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 公开路由
	r.GET("/login", views.LoginPage)
	r.POST("/login", views.Login)
	r.GET("/logout", views.Logout)

	authorized := r.Group("/")
	authorized.Use(AuthRequired())
	{
		authorized.GET("/", views.GetStudentList)
		authorized.POST("/update_password", views.UpdatePassword) // 修改密码接口

		// 学生操作
		studentGroup := authorized.Group("/student")
		{
			studentGroup.POST("/add", views.AddStudent)
			studentGroup.GET("/delete/:id", views.DeleteStudent)
			//studentGroup.GET("/edit/:id", views.EditStudentPage) // 稍后实现
		}

		// 4. 需要管理员权限的路由组
		adminGroup := authorized.Group("/admin")
		adminGroup.Use(AdminRequired()) // 挂载管理员拦截中间件
		{
			adminGroup.GET("/", views.AdminDashBoard)
			adminGroup.POST("/add_teacher", views.AddTeacher)
			adminGroup.GET("/delete_teacher/:id", views.DeleteTeacher)
			adminGroup.GET("/grant/:id", views.GrantAdmin)
			adminGroup.GET("/revoke/:id", views.RevokeAdmin)
		}
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get("userId")
		if userId == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		isAdmin := session.Get("isAdmin")
		if isAdmin == nil || isAdmin.(bool) == false {
			c.String(http.StatusForbidden, "403 Forbidden: 权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}
