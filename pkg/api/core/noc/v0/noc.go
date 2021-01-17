package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input noc.NOC

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, user.Result{Status: false, Error: err.Error()})
		return
	}

	if _, err := dbNOC.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, noc.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, noc.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, noc.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbNOC.Delete(&noc.NOC{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, noc.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, noc.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input noc.NOC

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, noc.Result{Status: false, Error: err.Error()})
		return
	}

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, user.Result{Status: false, Error: err.Error()})
		return
	}

	tmp := dbNOC.Get(noc.ID, &noc.NOC{Model: gorm.Model{ID: uint(id)}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, noc.Result{Status: false, Error: tmp.Err.Error()})
		return
	}

	if err := dbNOC.Update(noc.UpdateAll, replace(input, tmp.NOC[0])); err != nil {
		c.JSON(http.StatusInternalServerError, noc.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, noc.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, noc.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbNOC.Get(noc.ID, &noc.NOC{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, noc.Result{Status: false, Error: result.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, noc.Result{Status: true, NOC: result.NOC})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbNOC.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, noc.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, noc.Result{Status: true, NOC: result.NOC})
	}
}
