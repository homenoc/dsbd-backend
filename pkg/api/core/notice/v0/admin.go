package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbNotice "github.com/homenoc/dsbd-backend/pkg/api/store/notice/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input notice.Notice

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	log.Println(input.StartTime)

	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, token.Result{Status: false, Error: err.Error()})
		return
	}

	if _, err := dbNotice.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, notice.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, notice.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, notice.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbNotice.Delete(&notice.Notice{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, notice.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, notice.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input notice.Notice

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	tmp := dbNotice.Get(notice.ID, &notice.Notice{Model: gorm.Model{ID: input.ID}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, notice.Result{Status: false, Error: tmp.Err.Error()})
		return
	}

	if err := dbNotice.Update(notice.UpdateAll, updateAdminUser(input, tmp.Notice[0])); err != nil {
		c.JSON(http.StatusInternalServerError, notice.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, notice.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, notice.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbNotice.Get(notice.ID, &notice.Notice{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, notice.Result{Status: false, Error: result.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, notice.Result{Status: true, Notice: result.Notice})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbNotice.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, notice.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, notice.Result{Status: true, Notice: result.Notice})
	}
}
