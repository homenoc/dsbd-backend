package v0

import (
	"github.com/gin-gonic/gin"
	authInterface "github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnic"
	dbJPNICTech "github.com/homenoc/dsbd-backend/pkg/store/group/network/jpnicTech/v0"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input jpnic.Jpnic

	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: err.Error()})
		return
	}
	c.BindJSON(&input)

	if _, err := dbJPNICTech.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnic.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbJPNICTech.Delete(&jpnic.Jpnic{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnic.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input jpnic.Jpnic

	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: err.Error()})
		return
	}
	c.BindJSON(&input)

	if err := dbJPNICTech.Update(jpnic.UpdateAll, input); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnic.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbJPNICTech.Get(jpnic.ID, &jpnic.Jpnic{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnic.Result{Status: true, Jpnic: result.Jpnic})
}

func GetAllAdmin(c *gin.Context) {
	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: err.Error()})
		return
	}

	if result := dbJPNICTech.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnic.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, jpnic.Result{Status: true, Jpnic: result.Jpnic})
	}
}
