package worker

import (
	"net/http"
	"ta-kasir/base"
	"ta-kasir/config"
	"ta-kasir/helper"
	"ta-kasir/model"
	"ta-kasir/model/request"
	"ta-kasir/model/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AddWorker(c *gin.Context) {
	addWorker := request.AddWorker{}

	err := c.ShouldBindJSON(&addWorker)

	if err != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: err,
			Message: base.EmpetyField,
			Data: nil,
		})
		return
	}

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
	db := config.ConnectDatabase()

	var user model.User
	err = db.Debug().Where("email = ?", dataJWT.Email).
	Where("role = ?", 1).First(&user).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status: http.StatusInternalServerError,
			Error: err,
			Message: err.Error(),
			Data: nil,
		})
		return
	}

	hash, err  := bcrypt.GenerateFromPassword([]byte(addWorker.Password), bcrypt.DefaultCost)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: err,
			Message: err.Error(),
			Data: nil,
		})
		return
	}

worker := model.User{
	Username: addWorker.Username,
	Email: addWorker.Email,
	Password: string(hash),
	Role: 2,
}
	err = db.Debug().Create(&worker).Error

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
	Message: base.SuccessAddworker,
	Data: worker,
	})
}