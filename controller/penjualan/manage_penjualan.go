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

// type dataPesanan struct {
// 	NamaProduk string  `json:"nama_produk"`
// 	Quantity   int     `json:"quantity"`
// 	SubTotal   float64 `json:"sub_total"`
// }

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
		if err != nil {
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
	idPenjualan := fmt.Sprintf("TRS%s", formatTime+lastNumber)

	// Transaction Begin
	db.Transaction(func(tx *gorm.DB) error {
		// Insert tabel pelanggan
		err = tx.Debug().Create(&model.Pelanggan{
			IdPelanggan: idPelanggan,
			Email:       formAddPelanggan.Email,
			Nama:        formAddPelanggan.Nama,
			NoTelp:      formAddPelanggan.NoTelp,
			Alamat:      formAddPelanggan.Alamat,
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
			IdPenjualan:          idPenjualan,
			UserIduser:           dataJWT.UserId,
			PelangganIdPelanggan: idPelanggan,
			TotalHarga:           formAddPelanggan.Pembayaran.Grandtotal,
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
				IdDetailPenjualan:    uuid.New().String(),
				PenjualanIdPenjualan: idPenjualan,
				ProdukIdProduk:       each.IdProduk,
				JumlahProduk:         each.JumlahProduk,
				SubTotal:             each.SubTotal,
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
			Idpembayaran:         uuid.New().String(),
			PenjualanIdPenjualan: idPenjualan,
			Amount:               formAddPelanggan.Pembayaran.Amount,
			BiayaAdmin:           formAddPelanggan.Pembayaran.BiayaAdmin,
			Grandtotal:           formAddPelanggan.Pembayaran.Grandtotal,
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

		// get response transaksi

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

	// VERSI PNG
	imagePath, err := helper.GenerateImage(300, 400, idPenjualan, dataJWT.Nama, produkID, formAddPelanggan.DataPesanan, subTotals, formAddPelanggan.Pembayaran)

	// Versi PDF
	// imagePath, err := helper.GeneratePDF(idPenjualan, dataJWT.Nama, produkID, formAddPelanggan.DataPesanan, subTotals, formAddPelanggan.Pembayaran)

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
		Data: gin.H{
			"id_transaksi":      idPenjualan,
			"nama_kasir":        dataJWT.Nama,
			"tanggal_transaksi": time.Now().Format("02.01.2006"),
			"data_pesanan":      formAddPelanggan.DataPesanan,
			"amount":            formAddPelanggan.Pembayaran.Amount,
			"total_transaksi":   formAddPelanggan.Pembayaran.Grandtotal,
		},
	})
}

func ListTransaksi(c *gin.Context) {
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

	// PAGINATION
	limit := c.Query("limit")

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
				Status:  http.StatusBadRequest,
				Error:   err,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		db = db.Limit(limitInt)
	} else {
		db = db.Limit(10)
	}

	offset := c.Query("offset")

	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
				Status:  http.StatusBadRequest,
				Error:   err,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		db = db.Offset(offsetInt)
	} else {
		db = db.Offset(0)
	}

	// SEARCH
	key := c.Query("key")
	if key != "" {
		db = db.Where("id_penjualan LIKE ?", "%"+key+"%")
	}

	// FILTER BY
	layout := "2006-01-02"
	tanggalMulai := c.Query("tanggal_awal")
	tanggalAkhir := c.Query("tanggal_akhir")

	// fmt.Println(tanggalMulai)
	// fmt.Println(tanggalAkhir)

	if tanggalMulai != "" && tanggalAkhir != "" {
		_, err = time.Parse(layout, tanggalMulai)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
				Status:  http.StatusBadRequest,
				Error:   err,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		_, err = time.Parse(layout, tanggalAkhir)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
				Status:  http.StatusBadRequest,
				Error:   err,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		db = db.Where("penjualan.created_at >= ? AND penjualan.created_at <= ?", tanggalMulai, tanggalAkhir)
	}

	var transaksi []model.Penjualan
	err = db.Debug().Where("hapus = ?", 0).Order("created_at DESC").
		Preload("DetailPenjualan", func(tx *gorm.DB) *gorm.DB {
			return tx.Where("detail_penjualan.hapus = ?", 0)
		}).
		Preload("User", func(tx *gorm.DB) *gorm.DB {
			return tx.Where("user.hapus = ?", 0)
		}).
		Find(&transaksi).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSONP(http.StatusOK, response.Response{
		Status:  http.StatusOK,
		Error:   nil,
		Message: base.SuccessListTransaksi,
		Data:    transaksi,
	})
}

func DetailTransaksi(c *gin.Context) {
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

	idTransaksi := c.Query("idtransaksi")

	if idTransaksi == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   errors.New(base.ParamEmpty),
			Message: base.ParamEmpty + " idtransaksi kosong",
			Data:    nil,
		})
		return
	}

	db := config.ConnectDatabase()

	var detailPenjualan []model.DetailPenjualan

	err = db.Debug().
		Preload("Produk", func(tx *gorm.DB) *gorm.DB {
			return tx.Where("produk.hapus = ?", 0)
		}).
		Where("penjualan_id_penjualan = ?", idTransaksi).
		Where("hapus = ?", 0).
		Find(&detailPenjualan).Error

	length := len(detailPenjualan)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.ResponseArray{
		Status:  http.StatusOK,
		Error:   nil,
		Message: base.SuccessDetailTransaksi,
		Data:    detailPenjualan,
		Length:  length,
	})
}

func RefundTransaksi(c *gin.Context) {
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

	idDetailTransaksi := c.Query("iddetail")
	if idDetailTransaksi == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   errors.New(base.ParamEmpty + "(iddetail harus diisi)"),
			Message: base.ParamEmpty,
			Data:    nil,
		})
		return
	}

	db := config.ConnectDatabase()

	var detailPenjualan model.DetailPenjualan

	// query untuk mendapatkan detail transaksi berdasarkan id detail transaksi

	db.Where("id_detail_penjualan = ?", idDetailTransaksi).
		First(&detailPenjualan)

	idproduk := detailPenjualan.ProdukIdProduk
	idpenjualan := detailPenjualan.PenjualanIdPenjualan

	db.Transaction(func(tx *gorm.DB) error {
		// update produk stok
		err = tx.Debug().Model(&model.Produk{}).Where("id_produk = ?", idproduk).
			Update("stok", gorm.Expr("stok + ?", detailPenjualan.JumlahProduk)).Error
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

		// update pembayaran pada amount dan grandtotal
		err = tx.Debug().Model(&model.Pembayaran{}).Where("penjualan_id_penjualan = ?", idpenjualan).
			Update("amount", gorm.Expr("amount - ?", detailPenjualan.SubTotal)).
			Update("grandtotal", gorm.Expr("grandtotal - ?", detailPenjualan.SubTotal)).Error

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

		// update penjualan pada harga total di kurangi data yang di refund
		err = tx.Debug().Model(&model.Penjualan{}).Where("id_penjualan = ?", idpenjualan).
			Update("total_harga", gorm.Expr("total_harga - ?", detailPenjualan.SubTotal)).Error

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

		// update detail penjualan
		err = tx.Debug().Model(&model.DetailPenjualan{}).Where("id_detail_penjualan = ?", idDetailTransaksi).
			Update("hapus", 1).Error

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

		tx.Commit()
		return nil
	})

	c.JSON(http.StatusOK, response.Response{
		Status:  http.StatusOK,
		Error:   nil,
		Message: base.SuccessRefundTransaksi,
		Data:    nil,
	})
}

type ResponseListTransaksi struct {
	IdTransaksi      string    `json:"idtransaksi"`
	TanggalTransaksi time.Time `json:"tanggaltransaksi"`
	TotalPrice       float64   `json:"totalprice"`
	Totalitems       int       `json:"totalitems"`
}

func ListTransaksiV2(c *gin.Context) {
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
		db = db.Where("penjualan.id_penjualan LIKE ?", "%"+key+"%")
	}

	var transaksi []ResponseListTransaksi

	err = db.Debug().Table("penjualan").
		Select("penjualan.id_penjualan as idtransaksi, penjualan.created_at as tanggaltransaksi,"+
			"SUM(detail_penjualan.sub_total) as totalprice, COUNT(detail_penjualan.id_detail_penjualan) as totalitems").
		Joins("JOIN detail_penjualan ON detail_penjualan.penjualan_id_penjualan = detail_penjualan.penjualan_id_penjualan").
		Where("penjualan.hapus = ?", 0).
		Group("penjualan.id_penjualan").
		Find(&transaksi).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status:  http.StatusOK,
		Error:   nil,
		Message: "Success",
		Data:    transaksi,
	})
}

func GetTotalPendapatan(c *gin.Context) {
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

	// validasi admin
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

	dataPendapatan := response.ResponsePendapatan{}
	err = db.Table("pembayaran").Where("pembayaran.hapus = ?", 0).
		Select("SUM(grandtotal) as pendapatan_kotor," +
			"SUM(amount) as pendapatan_bersih").
		Find(&dataPendapatan).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status:  http.StatusOK,
		Error:   nil,
		Message: "Success",
		Data:    dataPendapatan,
	})

}
