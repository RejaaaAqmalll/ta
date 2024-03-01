package produk

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"ta-kasir/base"
	"ta-kasir/config"
	"ta-kasir/helper"
	"ta-kasir/model"
	"ta-kasir/model/request"
	"ta-kasir/model/response"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func ListProdukPetugas(c *gin.Context) {
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

	db := config.ConnectDatabase()

	key := c.Query("key")
	if key != "" {
		db = db.Where("nama_produk LIKE ?", "%"+key+"%")
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

func GetProdukByIdPetugas(c *gin.Context) {
	dataJWt, err := helper.GetClaims(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   err,
			Message: base.NoUserLogin,
			Data:    nil,
		})
		return
	}

	isPetugas := dataJWt.Role == 3

	if !isPetugas {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   errors.New(base.ShouldPetugas),
			Message: base.ShouldPetugas,
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

	var dataProduk model.Produk
	err = db.Debug().Where("id_produk = ?", idProduk).
		First(&dataProduk).Error

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
		Data:    dataProduk,
	})
}

func AddProdukPetugas(c *gin.Context) {
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

	formAddProduk := request.AddProduk{}

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.ParamEmpty,
			Data:    nil,
		})
		return
	}

	err = c.ShouldBind(&formAddProduk)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.ParamEmpty,
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
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status:  http.StatusOK,
		Error:   nil,
		Message: base.SuccessAddProduk,
		Data:    produk,
	})
}

func EditProdukPetugas(c *gin.Context) {
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

	// validasi akses hanya untuk petugas
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

	// parameter id untuk kondisi pengeditan data
	idProduk := c.Param("id")
	if idProduk == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.ParamEmpty,
			Data:    nil,
		})
		return
	}

	// Binding form dengan dengan id dan name "file"
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

	// request body dan Binding
	formEdit := request.EditProduk{}
	err = c.ShouldBind(&formEdit)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.EmpetyField,
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
		NamaProduk: formEdit.NamaProduk,
		Harga:      formEdit.Harga,
		Stok:       formEdit.Stok,
		Gambar:     file.Filename,
		LinkGambar: finalLink,
	}

	err = db.Debug().Model(model.Produk{}).
		Where("id_produk = ?", idProduk).
		Updates(&produk).Error

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
		Message: base.SuccessEditPorduk,
		Data: gin.H{
			"data_produk": produk,
			"link":        finalLink,
		},
	})
}

func DeleteProdukPetugas(c *gin.Context) {
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

	isPetugas := dataJWT.Role == 3

	if !isPetugas {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.ShouldPetugas,
			Data:    nil,
		})
		return
	}

	idProduk := c.Param("id")
	if idProduk == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.ParamEmpty,
			Data:    nil,
		})
		return
	}

	db := config.ConnectDatabase()

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
