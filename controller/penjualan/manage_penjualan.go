package penjualan

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"ta-kasir/base"
	"ta-kasir/config"
	"ta-kasir/helper"
	"ta-kasir/model"
	"ta-kasir/model/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Transaksi struct {
	Email string `json:"email" form:"email" binding:"email"`
	Nama string `json:"nama" form:"nama" binding:"required"`
	NoTelp string `json:"no_telp" form:"no_telp" binding:"required"`
	Alamat string `json:"alamat" form:"alamat" binding:"required"`
	DataPesanan []Pesanan `json:"data_pesanan" form:"data_pesanan" binding:"required"`
}

type Pesanan struct {
	IdProduk int `json:"id_produk" form:"id_produk" binding:"required"`
	JumlahProduk int `json:"jumlah_produk" form:"jumlah_produk" binding:"required"`
	TotalHarga int `json:"total_harga" form:"total_harga" binding:"required"`
}


func AddPenjualan(c *gin.Context) {
	dataJWT, err := helper.GetClaims(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   err,
			Message: base.NoUserLogin,
			Data:    nil,
		})
		return
	}

	// validasi role wajib 3
	isPetugas := dataJWT.Role == 3
	if !isPetugas {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   errors.New(base.ShouldPetugas),
			Message: base.ShouldPetugas,
			Data:    nil,
		})
		return
	}

	formAddPelanggan := Transaksi{}

	err = c.ShouldBind(&formAddPelanggan)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.EmpetyField,
			Data:    nil,
		})
		return
	}

	var pelanggan model.Pelanggan
	var lastID string
	db := config.ConnectDatabase()
	var produk model.Produk
	// validasi stok barang

	for _, each := range formAddPelanggan.DataPesanan {
		err = db.Debug().Where("id_produk = ?", each.IdProduk).Find(&produk).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
				Status:  http.StatusInternalServerError,
				Error:   err,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		fmt.Println(produk)

		if each.JumlahProduk > produk.Stok {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
				Status:  http.StatusBadRequest,
				Error:   err,
				Message: base.OutOfStock,
				Data:    nil,
			})
			return
		}
	}
	
	// mengambil id terakhir pelanggan
	err = db.Last(&pelanggan).Error

	if err == gorm.ErrRecordNotFound {
		if pelanggan.IdPelanggan == "" {
			lastID = "PLG000"
		} else {
			lastID = pelanggan.IdPelanggan
		}
	}


	lastNum, err := strconv.Atoi(lastID[3:])
	if err !=  nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	newNum := lastNum + 1

	idPelanggan := fmt.Sprintf("PLG%03d", newNum)
	db.Transaction(func(tx *gorm.DB) error {
		// Insert tabel pelanggan
		err = tx.Create(&model.Pelanggan{
			IdPelanggan: idPelanggan,
			Email:       formAddPelanggan.Email,
			Nama:        formAddPelanggan.Nama,
			NoTelp:      formAddPelanggan.NoTelp,
			Alamat: 	formAddPelanggan.Alamat,
		}).Error

		if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return err
		}


		// for _, each := range formAddPelanggan.DataPesanan {
			
		// }
		


		// ketika semua tidak error maka lakukan commit
		tx.Commit()
		return nil
	})

	
	
}