package main

import (
	"net/http"
	"os"
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
	route.StaticFS("/foto", gin.Dir("./storage/foto", false))
	route.POST("/register", auth.Register)
	route.POST("/login", auth.Login)

	route.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	

	err := os.Chmod("./storage/foto", 0755)

	if err != nil {
		panic(err)
	}
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
		admin.GET("/list_produk", produk.ListProduk)
		admin.GET("/get_produk/:id", produk.GetProdukById)
		admin.POST("/add_produk", produk.AddProduk)
		admin.PATCH("/edit_produk/:id", produk.EditProduk)
		admin.DELETE("/delete_produk/:id", produk.DeleteProduk)
	}

// ====================== PETUGAS ======================================

	petugas := route.Group("/petugas")
	{
		petugas.Use(authmiddleware())

		// Produk
		petugas.GET("/list_produk", produk.ListProdukPetugas)
		petugas.GET("/get_produk/:id", produk.GetProdukByIdPetugas)
		petugas.POST("/add_produk", produk.AddProdukPetugas)



		petugas.POST("/add_penjualan",)
	}


	route.Run(":8080")
}