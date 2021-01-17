package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	router "github.com/homenoc/dsbd-backend/pkg/api/core/noc/router"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbRouter "github.com/homenoc/dsbd-backend/pkg/api/store/noc/router/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input router.Router

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

	if _, err := dbRouter.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, router.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, router.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbRouter.Delete(&router.Router{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, router.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, router.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input router.Router

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

	tmp := dbRouter.Get(router.ID, &router.Router{Model: gorm.Model{ID: uint(id)}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, router.Result{Status: false, Error: tmp.Err.Error()})
		return
	}

	if err := dbRouter.Update(router.UpdateAll, replace(input, tmp.Router[0])); err != nil {
		c.JSON(http.StatusInternalServerError, router.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, router.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbRouter.Get(router.ID, &router.Router{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, router.Result{Status: false, Error: result.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, router.Result{Status: true, Router: result.Router})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbRouter.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, router.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, router.Result{Status: true, Router: result.Router})
	}
}
