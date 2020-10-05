package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/store/group/network/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/store/group/v0"
	"github.com/jinzhu/gorm"
	"net/http"
)

func Add(c *gin.Context) {
	var input network.Network
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	if err := check(input); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}

	if result.Group.Status == 500 {
		result.Group.Status = 100 + input.Type*10 + 1
	} else if result.Group.Status == 2 {
		result.Group.Status = input.Type*10 + 1
	} else {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: "error: group status"})
		return
	}

	_, err := dbNetwork.Create(&network.Network{
		GroupID: result.Group.ID, Type: input.Type, Name: input.Name, IP: input.IP, Route: input.Route,
		Date: input.Date, Plan: input.Plan})
	if err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
		return
	}
	if err := dbGroup.Update(group.UpdateStatus, group.Group{Model: gorm.Model{ID: result.Group.ID},
		Status: result.Group.Status + 1}); err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, group.Result{Status: true})
}

func Confirm(c *gin.Context) {
	var input network.Network
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: result.Err.Error()})
		return
	}
	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	if !(result.Group.Status%100 == 11 || result.Group.Status%100 == 21 || result.Group.Status%100 == 31) &&
		result.Group.Status < 200 {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: fmt.Sprint("error: status error")})
		return
	}

	if err := dbGroup.Update(group.UpdateStatus, group.Group{Model: gorm.Model{ID: result.Group.ID},
		Status: result.Group.Status + 1}); err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, group.Result{Status: true})
}
