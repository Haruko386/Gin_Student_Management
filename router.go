package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func initRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
}
