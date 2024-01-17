package produk

import (
	"net/http"
	"ta-kasir/base"
	"ta-kasir/config"
	"ta-kasir/model"
	"ta-kasir/model/request"
	"ta-kasir/model/response"

	"github.com/gin-gonic/gin"
)

func AddProduk(c *gin.Context) {
	// _, err := helper.GetClaims(c)
	
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
	// 		Status:  http.StatusUnauthorized,
	// 		Error:   err,
	// 		Message: base.NoUserLogin,
	// 		Data:    nil,
	// 	})
	// }

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

	err = c.ShouldBindJSON(&formAddProduk)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.EmpetyField,
			Data:    nil,
		})
		return
	}

	db := config.ConnectDatabase()

	var produk  = model.Produk{
	NamaProduk: formAddProduk.NamaProduk,
	Harga:      formAddProduk.Harga,
	Stok:       formAddProduk.Stok,
	Gambar: 	file.Filename,	
	}

	err = db.Debug().Create(&produk).Error

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
		Message: base.SuccessAddProduk,
		Data:    produk,
	})
}