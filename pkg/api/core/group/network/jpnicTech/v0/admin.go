package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicTech"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbJPNICTech "github.com/homenoc/dsbd-backend/pkg/store/group/network/jpnicTech/v0"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input jpnicTech.JpnicTech

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	if _, err := dbJPNICTech.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicTech.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnicTech.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jpnicTech.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbJPNICTech.Delete(&jpnicTech.JpnicTech{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicTech.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnicTech.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input jpnicTech.JpnicTech

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	if err := dbJPNICTech.Update(jpnicTech.UpdateAll, input); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicTech.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnicTech.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jpnicTech.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbJPNICTech.Get(jpnicTech.ID, &jpnicTech.JpnicTech{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnicTech.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnicTech.Result{Status: true, Jpnic: result.Jpnic})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbJPNICTech.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnicTech.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, jpnicTech.Result{Status: true, Jpnic: result.Jpnic})
	}
}
