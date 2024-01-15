package main

import (
	"net/http"
	"ta-kasir/config"
	"ta-kasir/controller/auth"
	"ta-kasir/controller/worker"
	"ta-kasir/middleware"

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

	authmiddleware := middleware.AuthCheck
	route.POST("/register", auth.Register)
	route.POST("/login", auth.Login)

	route.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

// ====================== ADMIN ======================================
	admin := route.Group("/admin")
	{
		admin.Use(authmiddleware())

		// Worker
		admin.POST("/add_worker", worker.AddWorker)
		admin.PATCH("/edit_worker/:id", worker.EditWorker)
		admin.DELETE("/delete_worker/:id", worker.DeleteWorker)
	}

// ====================== PETUGAS ======================================

	petugas := route.Group("/petugas")
	{
		petugas.Use(authmiddleware())
		petugas.POST("/add_penjualan",)
	}


route.Run(":8080")
}
// user.GET("/test", func(c *gin.Context) {
// dataJWT, err := helper.GetClaims(c)
// if err != nil {
// 	c.JSON(500, gin.H{
// 		"message": err,
// 	})
// 	return
// }

// nama := dataJWT.Nama
// 	c.JSON(200, gin.H{
// 		"message": nama,
// 	})
// })
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