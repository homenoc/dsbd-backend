package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbJpnicAdmin "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/jpnicAdmin/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input jpnicAdmin.JpnicAdmin

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	log.Println(c.BindJSON(&input))

	if _, err := dbJpnicAdmin.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnicAdmin.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, jpnicAdmin.Result{Status: false, Error: err.Error()})
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

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	log.Println(c.BindJSON(&input))

	if err := dbJpnicAdmin.Update(jpnicAdmin.UpdateAll, input); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnicAdmin.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbJpnicAdmin.Get(jpnicAdmin.ID, &jpnicAdmin.JpnicAdmin{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: result.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, jpnicAdmin.Result{Status: true, Jpnic: result.Jpnic})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbJpnicAdmin.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, jpnicAdmin.Result{Status: true, Jpnic: result.Jpnic})
	}
}
