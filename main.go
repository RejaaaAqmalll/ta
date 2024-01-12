package main

import (
	"net/http"
	"ta-kasir/config"
	"ta-kasir/controller/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	config.ConnectDatabase()
	route.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	route.POST("/register", auth.Register)
	route.POST("/login", auth.Login)

	user := route.Group("/user")
	{
		user.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "test",
			})
		})
	}
	
	route.Run(":8080")
}
// user.GET("/login", func(c *gin.Context) {
// 	c.JSON(200, gin.H{
// 		"message": "login",
// 	})
// })
// user.GET("/test", func(c *gin.Context) {
// 	db := config.ConnectDatabase()
// 	user := model.User{}
// 	db.Where("iduser = ?", 1).Find(&user)

// 	c.JSON(200, user)
// })