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
	"gorm.io/gorm"
)

func AddWorker(c *gin.Context) {
	addWorker := request.AddWorker{}

	err := c.ShouldBind(&addWorker)

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
	Where("role = ?", 1).Where("hapus = ?", 0).First(&user).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: err,
			Message: base.ShouldAdmin,
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

func EditWorker(c *gin.Context)  {
	idWorker := c.Param("id")
	if idWorker == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: nil,
			Message: base.ParamEmpty,
			Data: nil,
		})
		return
	}

	editWorker := request.EditWorker{}

	err := c.ShouldBind(&editWorker)

	if err != nil {
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
	Where("role = ?", 1).Where("hapus = ?", 0).First(&user).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: err,
			Message: base.ShouldAdmin,
			Data: nil,
		})
		return
	}

	worker := model.User{
		Username: editWorker.Username,
		Email: editWorker.Email,
	}

	err = db.Debug().Where("iduser = ?", idWorker).
	Updates(&worker).Where("role = ?", 2).Where("hapus = ?", 0).Error

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
		Message: base.SuccessEditWorker,
		Data: worker,
	})
}

func DeleteWorker(c *gin.Context) {
	idWorker := c.Param("id")
	if idWorker == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: nil,
			Message: base.ParamEmpty,
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
	Where("role = ?", 1).Where("hapus = ?", 0).First(&user).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: err,
			Message: base.ShouldAdmin,
			Data: nil,
		})
		return
	}
	
	err = db.Model(&model.User{}).
	Exec("Update user SET hapus = ? WHERE iduser = ?", 1, idWorker).Error

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
		Message: base.SuccessDeleteWorker,
		Data: nil,
	})
}

func ListWorker(c *gin.Context)  {
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

	key := c.Query("key")
	if key != "" {
		db = db.Where("email LIKE ? OR username LIKE ?", "%"+key+"%", "%"+key+"%")
	}

	db = db.Order("iduser ASC")
	
	dbAuth := config.ConnectDatabase()

	err = dbAuth.Debug().Where("email = ?", dataJWT.Email).
	Where("role = ?", 1).Where("hapus = ?", 0).First(&model.User{}).Error

	if  err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   err,
			Message: base.ShouldAdmin,
			Data:    nil,
		})
		return
	}

	var workers []model.User
	err = db.Debug().Where("role = ?", 2).Where("hapus = ?", 0).
	Find(&workers).Error

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
		Message: base.SuccessListWorker,
		Data: workers,
	})
}

func GetWorkerById(c *gin.Context)  {
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

	idWorker := c.Param("id")

	db := config.ConnectDatabase()

	err = db.Debug().Where("email = ?", dataJWT.Email).
	Where("role = ?", 1).Where("hapus = ?", 0).First(&model.User{}).Error

	if  err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   err,
			Message: base.ShouldAdmin,
			Data:    nil,
		})
		return
	}

	var worker model.User

	err = db.Debug().Where("iduser = ?", idWorker).
	Where("role = ?", 2).
	First(&worker).Error

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
		Message: base.SuccessGetWorker,
		Data: worker,
	})
}