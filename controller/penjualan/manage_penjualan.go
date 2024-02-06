package penjualan

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"ta-kasir/base"
	"ta-kasir/config"
	"ta-kasir/helper"
	"ta-kasir/model"
	"ta-kasir/model/request"
	"ta-kasir/model/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

	formAddPelanggan := request.Transaksi{}

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
	db := config.ConnectDatabase()
	var produk []model.Produk
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
				
		for _, eachProduk := range produk {
			if each.JumlahProduk > eachProduk.Stok {
				c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
					Status:  http.StatusBadRequest,
					Error:   errors.New(base.OutOfStock),
					Message: base.OutOfStock,
					Data:    nil,
				})
				return
			}
		}
		
	}
	
	// mengambil id terakhir pelanggan
	db.Last(&pelanggan)
	var lastID string
	var idPelanggan string
	if pelanggan.IdPelanggan == "" {
		lastID = "PLG000"
		} else {
			lastID = pelanggan.IdPelanggan
		}
		

		if len(lastID) >= 3 {
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
			idPelanggan = fmt.Sprintf("PLG%03d", newNum)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
				Status:  http.StatusInternalServerError,
				Error:   errors.New(base.InvalidID),
				Message: base.InvalidID,
				Data:    nil,
			})
			return
		}




	// idPelanggan := fmt.Sprintf("PLG%03d", newNum)
	formatTime := time.Now().Format("060102")
	codeLength := 4
	var modalcode = "1234567890"

	rand.Seed(time.Now().UnixNano())
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = modalcode[rand.Intn(len(modalcode))]
	}
	lastNumber := string(code)
	idPenjualan :=  fmt.Sprintf("TRS%s", formatTime+lastNumber)
	
	// Transaction Begin
	db.Transaction(func(tx *gorm.DB) error {
		// Insert tabel pelanggan
		err = tx.Debug().Create(&model.Pelanggan{
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

		// masukkan ke tabel penjualan
		err = tx.Debug().Create(&model.Penjualan{
			IdPenjualan: idPenjualan,
			PelangganIdPelanggan: idPelanggan,
			TotalHarga: formAddPelanggan.Pembayaran.Grandtotal,	
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

		for _, each := range formAddPelanggan.DataPesanan {
			// loop , insert ke tabel detail penjualan
			err = tx.Debug().Create(&model.DetailPenjualan{
				IdDetailPenjualan: uuid.New().String(),
				PenjualanIdPenjualan: idPenjualan,
				ProdukIdProduk: each.IdProduk,
				JumlahProduk: each.JumlahProduk,
				SubTotal: each.SubTotal,
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
		}


		// insert tabel pembayaran
		err = tx.Debug().Create(&model.Pembayaran{
			Idpembayaran: uuid.New().String(),
			PenjualanIdPenjualan: idPenjualan,
			Amount: formAddPelanggan.Pembayaran.Amount,
			BiayaAdmin: formAddPelanggan.Pembayaran.BiayaAdmin,
			Grandtotal: formAddPelanggan.Pembayaran.Grandtotal,
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


		// update stok
		for _, each := range formAddPelanggan.DataPesanan {
			err = tx.Debug().Model(&model.Produk{}).Where("id_produk = ?", each.IdProduk).
			Update("stok", gorm.Expr("stok - ?", each.JumlahProduk)).Error
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
		}
		// ketika semua tidak error maka lakukan commit
		tx.Commit()
		return nil
		// akhir transaction
	})

	// Go Func untuk generate image dan send to email
	// go func() {

		var produkID []int
		var subTotals []float64

		for _, pesanan := range formAddPelanggan.DataPesanan {
			produkID = append(produkID, pesanan.IdProduk)
			subTotals = append(subTotals, pesanan.SubTotal)
		}


		imagePath, err := helper.GenerateImage(300, 400, idPenjualan, dataJWT.Nama, produkID, formAddPelanggan.DataPesanan, subTotals, formAddPelanggan.Pembayaran)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
				Status:  http.StatusInternalServerError,
				Error:   err,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		err = helper.SendEmail(formAddPelanggan.Email, imagePath)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
				Status:  http.StatusInternalServerError,
				Error:   err,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

	// }()

	c.JSON(http.StatusOK, response.Response{
		Status:  http.StatusOK,
		Error:   nil,
		Message: base.SuccesTransaction,
		Data: formAddPelanggan,
	})
}

func ListTransaksi(c *gin.Context) {
	dataJWT, err := helper.GetClaims(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status: http.StatusUnauthorized,
			Error:  err,
			Message: base.NoUserLogin,
			Data:   nil,
		})
		return
	}
	isAdmin := dataJWT.Role == 1 || dataJWT.Role == 2

	if !isAdmin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   errors.New(base.ShouldAdmin),
			Message: base.ShouldAdmin,
			Data:    nil,
		})
		return
	}

	db := config.ConnectDatabase()

	key := c.Query("key")
	if key != "" {
		db = db.Where("id_penjualan LIKE ?", "%"+key+"%")
	}

	var transaksi []model.Penjualan
	err = db.Debug().Where("hapus = ?", 0).Order("created_at DESC").Preload("DetailPenjualan").
	Find(&transaksi).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status: http.StatusInternalServerError,
			Error:  err,
			Message: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSONP(http.StatusOK, response.Response{
		Status: http.StatusOK,
		Error:  nil,
		Message: base.SuccessListTransaksi,
		Data:   transaksi,
	})
}