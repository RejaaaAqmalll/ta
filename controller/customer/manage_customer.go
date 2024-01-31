package customer

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

func ListCustomer(c *gin.Context) {
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
		db = db.Where("nama LIKE ? OR email LIKE ? OR no_telp LIKE ?", "%"+key+"%", "%"+key+"%", "%"+key+"%")
	}

	var listCustomer []model.Pelanggan

	err = db.Debug().
		Where("hapus = ?", 0).Order("id_pelanggan ASC").
		Find(&listCustomer).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status: http.StatusInternalServerError,
			Error:  err,
			Message: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status: http.StatusOK,
		Error:  nil,
		Message: base.SuccessListCustomer,
		Data:   listCustomer,
	})
}