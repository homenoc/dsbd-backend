package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/store/group/v0"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input group.Group

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	if _, err := dbGroup.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbGroup.Delete(&group.Group{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input group.Group

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	tmp := dbGroup.Get(group.ID, &group.Group{Model: gorm.Model{ID: input.ID}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: tmp.Err.Error()})
		return
	}

	replace, err := updateAdminUser(input, tmp.Group[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "error: this email is already registered"})
		return
	}

	if err := dbGroup.Update(group.UpdateAll, replace); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbGroup.Get(group.ID, &group.Group{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.Result{Status: true, Group: result.Group})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbGroup.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, group.Result{Status: true, Group: result.Group})
	}
}
