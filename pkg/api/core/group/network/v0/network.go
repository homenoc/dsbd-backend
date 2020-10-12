package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnic_user"
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
		Status: result.Group.Status}); err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, group.Result{Status: true})
}

func Update(c *gin.Context) {
	var input network.NetworkUser
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

	if err := checkUpdate(input); err != nil {
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

	resultNetwork := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: input.ID}})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: resultNetwork.Err.Error()})
		return
	}
	if len(resultNetwork.Network) == 0 {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: "failed Network ID"})
		return
	}
	if resultNetwork.Network[0].GroupID != result.Group.ID {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: "Authentication failure"})
		return
	}
	if resultNetwork.Network[0].Lock {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: "this network is locked..."})
		return
	}

	data := resultNetwork.Network[0]

	if input.Name != "" {
		data.Name = input.Name
	}
	if input.IP != "" {
		data.IP = input.IP
	}
	if input.Route != "" {
		data.Route = input.Route
	}
	if input.Date != "" {
		data.Date = input.Date
	}
	if input.Plan != "" {
		data.Plan = input.Plan
	}

	if err := dbNetwork.Update(network.UpdateData, data); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
	}

	c.JSON(http.StatusOK, group.Result{Status: true})
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
	c.JSON(http.StatusOK, group.Result{Status: true})
}
