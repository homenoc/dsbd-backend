package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	"github.com/jinzhu/gorm"
	"net/http"
)

func Add(c *gin.Context) {
	var input connection.Connection
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, connection.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusUnauthorized, connection.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	if !((result.Group.Status%100 == 13 || result.Group.Status%100 == 23) && (result.Group.Status/100 == 0 ||
		result.Group.Status/100 == 1)) {
		c.JSON(http.StatusUnauthorized, connection.Result{Status: false, Error: fmt.Sprint("error: status error")})
		return
	}

	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, connection.Result{Status: false, Error: err.Error()})
		return
	}

	_, err := dbConnection.Create(&connection.Connection{
		GroupID: result.Group.ID, UserID: input.UserID, Service: input.Service, NTT: input.NTT, NOC: input.NOC,
		TermIP: input.TermIP, Monitor: input.Monitor})
	if err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbGroup.Update(group.UpdateStatus, group.Group{Model: gorm.Model{ID: result.Group.ID},
		Status: result.Group.Status + 1}); err != nil {
		c.JSON(http.StatusInternalServerError, connection.Result{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, group.Result{Status: true})
}
