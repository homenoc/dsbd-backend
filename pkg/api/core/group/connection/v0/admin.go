package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input connection.Connection

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	log.Println(c.BindJSON(&input))

	if _, err := dbConnection.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, connection.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, connection.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbConnection.Delete(&connection.Connection{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, connection.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input connection.Connection

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	log.Println(c.BindJSON(&input))

	resultConnection := dbConnection.Get(connection.ID, &connection.Connection{Model: gorm.Model{ID: input.ID}})
	if resultConnection.Err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: resultConnection.Err.Error()})
		return
	}

	replace, err := updateAdminConnection(input, resultConnection.Connection[0])
	if err != nil {
		c.JSON(http.StatusUnauthorized, connection.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbConnection.Update(connection.UpdateAll, replace); err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, connection.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, connection.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbConnection.Get(connection.ID, &connection.Connection{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: result.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, connection.Result{Status: true, ConnectionData: result.Connection})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbConnection.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, connection.Result{Status: true, ConnectionData: result.Connection})
	}
}
