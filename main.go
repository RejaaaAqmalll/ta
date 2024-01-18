package main

import (
	"net/http"
	"ta-kasir/config"
	"ta-kasir/controller/auth"
	"ta-kasir/controller/produk"
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
		admin.GET("/list_worker", worker.ListWorker)
		admin.GET("/get_worker/:id", worker.GetWorkerById)
		admin.POST("/add_worker", worker.AddWorker)
		admin.PATCH("/edit_worker/:id", worker.EditWorker)
		admin.DELETE("/delete_worker/:id", worker.DeleteWorker)


		// Produk
		admin.POST("/add_produk", produk.AddProduk)
		admin.PATCH("/edit_produk/:id", produk.EditProduk)
		admin.DELETE("/delete_produk/:id", produk.DeleteProduk)
	}

// ====================== PETUGAS ======================================

	petugas := route.Group("/petugas")
	{
		petugas.Use(authmiddleware())
		petugas.POST("/add_penjualan",)
	}


	route.Run(":8080")
}