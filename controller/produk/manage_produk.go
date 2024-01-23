package produk

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"ta-kasir/base"
	"ta-kasir/config"
	"ta-kasir/helper"
	"ta-kasir/model"
	"ta-kasir/model/request"
	"ta-kasir/model/response"

	"github.com/gin-gonic/gin"
)

func AddProduk(c *gin.Context) {
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
	if err != nil{
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
	finalLink := "http://127.0.0.1:8080" + link

	var produk  = model.Produk{
	NamaProduk: formAddProduk.NamaProduk,
	Harga:      formAddProduk.Harga,
	Stok:       formAddProduk.Stok,
	Gambar: 	file.Filename,	
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
		Data:   produk,
	})
}

func EditProduk(c *gin.Context)  {
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
			Status: http.StatusBadRequest,
			Error: nil,
			Message: base.ParamEmpty,
			Data: nil,
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
	if err != nil{
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
	finalLink := "http://127.0.0.1:8080" + link

	var produk  = model.Produk{
		NamaProduk: formEditProduk.NamaProduk,
		Harga: formEditProduk.Harga,
		Stok: formEditProduk.Stok,
		Gambar: file.Filename,
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
		Data:    gin.H{
			"data_produk": produk,
			"link":        finalLink,
		},
	})
}

func DeleteProduk(c *gin.Context)  {
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
		Message: base.SuccessDeleteProduk,
		Data: nil,
	})
}

func ListProduk(c *gin.Context)  {
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

func GetProdukById(c *gin.Context)  {
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
		Message: base.SuccessGetProduk,
		Data: produk,
	})
}