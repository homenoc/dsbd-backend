package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	networkUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/network_user"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbNetworkUser "github.com/homenoc/dsbd-backend/pkg/store/group/network/network_user/v0"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/store/group/network/v0"
	"github.com/jinzhu/gorm"
	"net/http"
)

func Add(c *gin.Context) {
	var input networkUser.NetworkUser
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	if !(result.Group.Status%10 == 1 && result.Group.Status/100 >= 0 && result.Group.Status/100 <= 3) {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: "status error"})
		return
	}

	if err := checkGroupID(result.Group.ID, input.NetworkID, input.JPNICUserID); err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: err.Error()})
		return
	}

	if err := checkDuplicate(input); err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: err.Error()})
		return
	}

	_, err := dbNetworkUser.Create(&networkUser.NetworkUser{Type: input.Type, NetworkID: input.NetworkID, JPNICUserID: input.JPNICUserID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, networkUser.Result{Status: true})
}

func Delete(c *gin.Context) {
	var input networkUser.NetworkUser
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	if !(result.Group.Status%10 == 1 && result.Group.Status/100 >= 0 && result.Group.Status/100 <= 3) {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: "status error"})
		return
	}

	resultNetworkUser := dbNetworkUser.Get(networkUser.ID, &networkUser.NetworkUser{Model: gorm.Model{ID: input.ID}})
	if resultNetworkUser.Err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: resultNetworkUser.Err.Error()})
		return
	}
	if len(resultNetworkUser.NetworkUser) == 0 {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: "failed Network ID"})
		return
	}
	if !resultNetworkUser.NetworkUser[0].Lock {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: "this data is locked..."})
		return
	}

	resultNetwork := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: resultNetworkUser.NetworkUser[0].NetworkID}})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: resultNetwork.Err.Error()})
		return
	}
	if len(resultNetwork.Network) == 0 {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: "failed Network ID"})
		return
	}
	if resultNetwork.Network[0].GroupID != result.Group.ID {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: "not match GroupID"})
		return
	}

	if err := dbNetworkUser.Delete(&networkUser.NetworkUser{Model: gorm.Model{ID: input.ID}}); err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: err.Error()})
	}
	c.JSON(http.StatusOK, networkUser.Result{Status: true})
}

func Update(c *gin.Context) {
	var input networkUser.NetworkUser
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	if !(result.Group.Status%10 == 1 && result.Group.Status/100 >= 0 && result.Group.Status/100 <= 3) {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: "status error"})
		return
	}

	resultNetworkUser := dbNetworkUser.Get(networkUser.ID, &networkUser.NetworkUser{Model: gorm.Model{ID: input.ID}})
	if resultNetworkUser.Err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: resultNetworkUser.Err.Error()})
		return
	}
	if len(resultNetworkUser.NetworkUser) == 0 {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: "failed Network ID"})
		return
	}

	if err := checkGroupID(result.Group.ID, input.NetworkID, input.JPNICUserID); err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: err.Error()})
	}

	if err := checkDuplicate(input); err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbNetworkUser.Update(networkUser.UpdateInfo, input); err != nil {
		c.JSON(http.StatusInternalServerError, networkUser.Result{Status: false, Error: err.Error()})
	}

	c.JSON(http.StatusOK, networkUser.Result{Status: true})
}
