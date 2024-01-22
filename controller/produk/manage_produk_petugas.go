package produk

import (
	"errors"
	"net/http"
	"ta-kasir/base"
	"ta-kasir/config"
	"ta-kasir/helper"
	"ta-kasir/model"
	"ta-kasir/model/response"

	"github.com/gin-gonic/gin"
)

func ListProdukPetugas(c *gin.Context) {
	dataJWT, err := helper.GetClaims(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status: http.StatusUnauthorized,
			Error: err,
			Message: base.NoUserLogin,
			Data: nil,
		})
		return
	}

		isPetugas := dataJWT.Role == 2

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
			Status: http.StatusInternalServerError,
			Error: err,
			Message: err.Error(),
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status: http.StatusOK,
		Error:  nil,
		Message: base.SuccessListProduk,
		Data: listProduk,
	})
}

func GetProdukByIdPetugas(c *gin.Context) {
	dataJWt, err := helper.GetClaims(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status: http.StatusUnauthorized,
			Error: err,
			Message: base.NoUserLogin,
			Data: nil,
		})
		return
	}

	isPetugas := dataJWt.Role == 2

	if !isPetugas {
	    c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
	        Status:  http.StatusUnauthorized,
	        Error:   errors.New(base.ShouldAdmin),
	        Message: base.ShouldAdmin,
	        Data:    nil,
	    })
	    return
	}

	idProduk := c.Param("id")
	if idProduk == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: nil,
			Message: base.ParamEmpty,
			Data: nil,
		})
		return
	}

	db  := config.ConnectDatabase()

	var dataProduk model.Produk
	err = db.Debug().Where("id_produk = ?", idProduk).
	First(&dataProduk).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status: http.StatusInternalServerError,
			Error: err,
			Message: err.Error(),
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status: http.StatusOK,
		Error:  nil,
		Message: base.SuccessListProduk,
		Data: dataProduk,
	})
}

func AddProdukPetugas(c *gin.Context)  {
	dataJWT, err  := helper.GetClaims(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status: http.StatusUnauthorized,
			Error: err,
			Message: base.NoUserLogin,
			Data: nil,
		})
		return
	}


	isPetugas := dataJWT.Role == 2
	if !isPetugas {
	    c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
	        Status:  http.StatusUnauthorized,
	        Error:   errors.New(base.ShouldPetugas),
	        Message: base.ShouldPetugas,
	        Data:    nil,
	    })
	    return
	}


}