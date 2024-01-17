package auth

import (
	"net/http"
	base "ta-kasir/base"
	"ta-kasir/config"
	"ta-kasir/helper"
	"ta-kasir/model"
	"ta-kasir/model/request"
	"ta-kasir/model/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	formRegister := request.Register{}

	// Bind input
	err := c.ShouldBindJSON(&formRegister)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.EmpetyField,
			Data:    nil,
		})
		return
	}

	db := config.ConnectDatabase()

err = db.Debug().
Where("email = ?", formRegister.Email).Where("hapus = ?", 0).
First(&model.User{}).Error

if err != gorm.ErrRecordNotFound {
	c.JSON(http.StatusBadRequest, response.Response{
		Status: http.StatusBadRequest,
		Error: err,
		Message: base.AlreadyRegister,
		Data: nil,
	})
	return
}

hash, err := bcrypt.GenerateFromPassword([]byte(formRegister.ConfirmPassword), bcrypt.DefaultCost)

if err != nil {
	c.JSON(http.StatusBadRequest, response.Response{
		Status: http.StatusBadRequest,
		Error: err,
		Message: base.FailedhashPw,
		Data: nil,
	})
	return
}

user := model.User{
	Username: formRegister.Username,
	Email: formRegister.Email,
	Role: 1,
	Password: string(hash),
}
err = db.Debug().Create(&user).Error

if err != nil {
	c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
		Status: http.StatusInternalServerError,
		Error: err,
		Message: base.FailedCreateUser,
		Data: nil,
	})
	return
}

c.JSON(http.StatusOK, response.Response{
	Status: http.StatusOK,
	Error: nil,
	Message: base.SuccessRegister,
	Data: user,
})
}

func Login(c *gin.Context)  {
	formLogin := request.Login{}

	err := c.ShouldBindJSON(&formLogin)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: err,
			Message: base.EmpetyField,
			Data: nil,
		})
		return
	}

	db := config.ConnectDatabase()

	var user model.User
	err = db.Debug().Where("email = ?", formLogin.Email).
	Where("hapus = ?", 0).
	First(&user).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: err,
			Message: base.UserNotFound,
			Data: nil,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(formLogin.Password))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: err,
			Message: base.IncorrectPassEmail,
			Data: nil,
		})
		return
	}

	token, err := helper.GenerateToken(user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: err,
			Message: base.FailedGenerateToken,
			Data: nil,
		})
	}

	c.JSON(http.StatusOK, response.Response{
		Status: http.StatusOK,
		Error: nil,
		Message: base.SuccessLogin,
		Data: gin.H{
			"token": token,
			"data_user": user,
		},
	})
}