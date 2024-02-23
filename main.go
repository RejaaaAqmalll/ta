package main

import (
	"net/http"
	"ta-kasir/config"
	"ta-kasir/controller/auth"
	"ta-kasir/controller/customer"
	"ta-kasir/controller/penjualan"
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

		// Customer
		admin.GET("/list_customer", customer.ListCustomer)
		admin.DELETE("/delete_customer/:id", customer.DeleteCustomer)

		// Transaksi
		admin.GET("/list_transaksi", penjualan.ListTransaksi)
		admin.GET("/detail_transaksi", penjualan.DetailTransaksi)
		admin.PATCH("/edit_transaksi", penjualan.EditTransaksi)
	}

	// ====================== PETUGAS ======================================

	petugas := route.Group("/petugas")
	{
		petugas.Use(authmiddleware())

		// Produk
		petugas.GET("/list_produk", produk.ListProdukPetugas)
		petugas.GET("/get_produk/:id", produk.GetProdukByIdPetugas)
		petugas.POST("/add_produk", produk.AddProdukPetugas)
		petugas.PATCH("/edit_produk/:id", produk.EditProdukPetugas)
		petugas.DELETE("/delete_produk/:id", produk.DeleteProdukPetugas)

		petugas.POST("/add_penjualan", penjualan.AddPenjualan)
	}

	route.Run(":8080")
}
