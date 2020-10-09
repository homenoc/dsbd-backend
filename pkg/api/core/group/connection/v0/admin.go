package v0

import (
	"github.com/gin-gonic/gin"
	authInterface "github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/store/group/connection/v0"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input connection.Connection

	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}
	c.BindJSON(&input)

	if _, err := dbConnection.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, connection.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
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

	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}
	c.BindJSON(&input)

	if err := dbConnection.Update(connection.UpdateAll, input); err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, connection.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbConnection.Get(connection.ID, &connection.Connection{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, connection.Result{Status: true, ConnectionData: result.Connection})
}

func GetAllAdmin(c *gin.Context) {
	if err := auth.AdminAuthentication(authInterface.AdminStruct{User: c.Request.Header.Get("USER"),
		Pass: c.Request.Header.Get("PASS")}); err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}

	if result := dbConnection.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, connection.Result{Status: true, ConnectionData: result.Connection})
	}
}
