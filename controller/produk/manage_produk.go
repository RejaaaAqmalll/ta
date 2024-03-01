package produk

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"ta-kasir/base"
	"ta-kasir/config"
	"ta-kasir/helper"
	"ta-kasir/model"
	"ta-kasir/model/request"
	"ta-kasir/model/response"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AddProduk(c *gin.Context) {
	godotenv.Load()
	dataJWT, err := helper.GetClaims(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   err,
			Message: base.NoUserLogin,
			Data:    nil,
		})
	}

	formAddProduk := request.AddProduk{}

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	err = c.ShouldBind(&formAddProduk)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.EmpetyField,
			Data:    nil,
		})
		return
	}

	// validasi role wajin 1 atau 2

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

	// validasi input file harus berupa gambar
	src, err := file.Open()
	if err != nil {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	buffer := make([]byte, 261)
	_, err = src.Read(buffer)

	if err != nil && err != io.EOF {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// get mime type
	kind := http.DetectContentType(buffer)
	if kind == "" || !helper.IsSupportedImageFormat(kind) {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.FileNotSupported,
			Data:    nil,
		})
		return
	}

	fileName := helper.GenerateFilename(file.Filename)

	fileDest := helper.GetImageSavePath(fileName)

	err = c.SaveUploadedFile(file, fileDest)
	if err != nil {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	db := config.ConnectDatabase()

	link := fmt.Sprintf("/foto/%s", fileName)
	finalLink := os.Getenv("BASE_URL") + link

	var produk = model.Produk{
		NamaProduk: formAddProduk.NamaProduk,
		Harga:      formAddProduk.Harga,
		Stok:       formAddProduk.Stok,
		Gambar:     file.Filename,
		LinkGambar: finalLink,
	}

	err = db.Debug().Create(&produk).Error

	if err != nil {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// fmt.Println(link)
	// fmt.Println(finalLink)
	c.JSON(http.StatusOK, response.Response{
		Status:  http.StatusOK,
		Error:   nil,
		Message: base.SuccessAddProduk,
		Data:    produk,
	})
}

func EditProduk(c *gin.Context) {
	godotenv.Load()
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

	idProduk := c.Param("id")
	if idProduk == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   nil,
			Message: base.ParamEmpty,
			Data:    nil,
		})
		return
	}

	formEditProduk := request.EditProduk{}
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.ShouldBind(&formEditProduk)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.EmpetyField,
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

	// validasi input file harus berupa gambar
	src, err := file.Open()
	if err != nil {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	defer src.Close()

	buffer := make([]byte, 261)
	_, err = src.Read(buffer)

	if err != nil && err != io.EOF {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// get mime type
	kind := http.DetectContentType(buffer)
	if kind == "" || !helper.IsSupportedImageFormat(kind) {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.FileNotSupported,
			Data:    nil,
		})
		return
	}

	fileName := helper.GenerateFilename(file.Filename)

	fileDest := helper.GetImageSavePath(fileName)

	err = c.SaveUploadedFile(file, fileDest)
	if err != nil {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	db := config.ConnectDatabase()

	link := fmt.Sprintf("/foto/%s", fileName)
	finalLink := os.Getenv("BASE_URL") + link

	var produk = model.Produk{
		NamaProduk: formEditProduk.NamaProduk,
		Harga:      formEditProduk.Harga,
		Stok:       formEditProduk.Stok,
		Gambar:     file.Filename,
		LinkGambar: finalLink,
	}

	err = db.Debug().Model(model.Produk{}).
		Where("id_produk = ?", idProduk).
		Updates(&produk).Error

	if err != nil {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// fmt.Println(finalLink)
	c.JSON(http.StatusOK, response.Response{
		Status:  http.StatusOK,
		Error:   nil,
		Message: base.SuccessEditPorduk,
		Data: gin.H{
			"data_produk": produk,
			"link":        finalLink,
		},
	})
}

func DeleteProduk(c *gin.Context) {
	idProduk := c.Param("id")

	if idProduk == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   nil,
			Message: base.ParamEmpty,
			Data:    nil,
		})
		return
	}

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

	db := config.ConnectDatabase()
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

	err = db.Debug().Model(model.Produk{}).
		Where("id_produk = ?", idProduk).Update("hapus", 1).Error

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
		Message: base.SuccessDeleteProduk,
		Data:    nil,
	})
}

func ListProduk(c *gin.Context) {
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

	db := config.ConnectDatabase()

	key := c.Query("key")
	if key != "" {
		db = db.Where("nama_produk LIKE ?", "%"+key+"%")
	}

	// validasi admin terlebih dahulu
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

	var listProduk []model.Produk
	err = db.Debug().
		Where("hapus = ?", 0).Order("id_produk ASC").Find(&listProduk).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	var totalStok int64
	err = db.Debug().
		Table("produk").Select("SUM(stok) AS total_stok").Where("hapus = ?", 0).
		Find(&totalStok).Error
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
		Message: base.SuccessListProduk,
		Data: gin.H{
			"listProduk": listProduk,
			"total_stok": totalStok,
		},
	})
}

func GetProdukById(c *gin.Context) {
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

	idProduk := c.Param("id")
	if idProduk == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   nil,
			Message: base.ParamEmpty,
			Data:    nil,
		})
		return
	}

	db := config.ConnectDatabase()

	// validasi admin terlebih dahulu
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

	var produk model.Produk

	err = db.Debug().Where("id_produk = ?", idProduk).
		Where("hapus = ?", 0).First(&produk).Error

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
		Message: base.SuccessGetProduk,
		Data:    produk,
	})
}

type dataBestSeller struct {
	IdProduk    int     `json:"id_produk"`
	NamaProduk  string  `json:"nama_produk"`
	HargaProduk float64 `json:"harga_produk"`
	TotalHarga  float64 `json:"total_harga"`
	Terjual     int     `json:"terjual"`
}

func GetProdukBestSeller(c *gin.Context) {
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

	// validasi admin terlebih dahulu
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
	idprodukParam := c.Query("id_produk")

	// memisahkan id
	idprodukArray := strings.Split(idprodukParam, ",")

	var idprodukInt []int
	for _, idStr := range idprodukArray {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter"})
			return
		}
		idprodukInt = append(idprodukInt, id)
	}

	var listProduk []dataBestSeller

	if len(idprodukInt) > 0 {
		if err = db.Debug().Table("detail_penjualan").
			Select("SUM(detail_penjualan.jumlah_produk) as terjual,"+
				"detail_penjualan.produk_id_produk as id_produk,"+
				"produk.nama_produk as nama_produk,"+
				"produk.harga as harga_produk,"+
				"SUM(detail_penjualan.sub_total) as total_harga").
			Joins("JOIN produk ON produk.id_produk = detail_penjualan.produk_id_produk").
			Where("detail_penjualan.produk_id_produk IN (?)", idprodukInt).
			Where("detail_penjualan.hapus = ?", 0).
			Group("detail_penjualan.produk_id_produk").
			Order("detail_penjualan.produk_id_produk ASC").
			Find(&listProduk).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
				Status:  http.StatusInternalServerError,
				Error:   err,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
	}

	// get total terjual
	dbtotal := config.ConnectDatabase()

	var totalTerjual int
	err = dbtotal.Table("detail_penjualan").
		Select("SUM(detail_penjualan.jumlah_produk)").
		Find(&totalTerjual).Error

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
		Message: base.SuccessGetBestSeller,
		Data: gin.H{
			"data_best_seller": listProduk,
			"total_terjual":    totalTerjual,
		},
	})
}
