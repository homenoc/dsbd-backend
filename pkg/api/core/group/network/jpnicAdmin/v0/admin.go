package v0

import (
	"github.com/gin-gonic/gin"
	authInterface "github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	dbJpnicAdmin "github.com/homenoc/dsbd-backend/pkg/store/group/network/jpnicAdmin/v0"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input jpnicAdmin.JpnicAdmin

	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}
	c.BindJSON(&input)

	if _, err := dbJpnicAdmin.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnicAdmin.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbJpnicAdmin.Delete(&jpnicAdmin.JpnicAdmin{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnicAdmin.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input jpnicAdmin.JpnicAdmin

	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}
	c.BindJSON(&input)

	if err := dbJpnicAdmin.Update(jpnicAdmin.UpdateAll, input); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnicAdmin.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbJpnicAdmin.Get(jpnicAdmin.ID, &jpnicAdmin.JpnicAdmin{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnicAdmin.Result{Status: true, Jpnic: result.Jpnic})
}

func GetAllAdmin(c *gin.Context) {
	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}

	if result := dbJpnicAdmin.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, jpnicAdmin.Result{Status: true, Jpnic: result.Jpnic})
	}
}