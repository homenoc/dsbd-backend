package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
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

	//if result.Group.Status == 111 || 121 || 211 || 221 {
	if result.Group.Status == 2 {
		if input.PI {
			result.Group.Status = 21
		} else {
			result.Group.Status = 11
		}
		if err := dbGroup.Update(group.UpdateStatus, group.Group{Model: gorm.Model{ID: result.Group.ID},
			Status: result.Group.Status}); err != nil {
			c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
			return
		}
	} else if !(result.Group.Status == 111 || result.Group.Status == 121) {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: "error: group status"})
		return
	}

	net, err := dbNetwork.Create(&network.Network{
		GroupID: result.Group.ID, Org: input.Org, OrgEn: input.OrgEn, Postcode: input.Postcode, Address: input.Address,
		AddressEn: input.AddressEn, Route: input.Route, PI: input.PI, ASN: input.ASN, V4: input.V4, V6: input.V6,
		V4Name: input.V4Name, V6Name: input.V6Name, Date: input.Date, Plan: input.Plan, Lock: false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, network.ResultOne{Status: true, Network: *net})
}

func Update(c *gin.Context) {
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

	if !(result.Group.Status == 211 || result.Group.Status == 221) {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: "error: group status"})
		return
	}

	resultNetwork := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: input.ID}})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: resultNetwork.Err.Error()})
		return
	}
	if len(resultNetwork.Network) == 0 {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: "failed Network ID"})
		return
	}
	if resultNetwork.Network[0].GroupID != result.Group.ID {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: "Authentication failure"})
		return
	}
	if resultNetwork.Network[0].Lock {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: "this network is locked..."})
		return
	}

	replace := replaceNetwork(resultNetwork.Network[0], input)

	if err := dbNetwork.Update(network.UpdateData, replace); err != nil {
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

	if !((result.Group.Status%100 == 11 || result.Group.Status%100 == 21) && (result.Group.Status/100 == 0 ||
		result.Group.Status/100 == 1)) {
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
