package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	jpnicUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnic_user"
	networkUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/network_user"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/store/group/connection/v0"
	dbJPNICUser "github.com/homenoc/dsbd-backend/pkg/store/group/network/jpnic_user/v0"
	dbNetworkUser "github.com/homenoc/dsbd-backend/pkg/store/group/network/network_user/v0"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/store/group/network/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/store/group/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	"github.com/jinzhu/gorm"
	"net/http"
)

//参照関連のエラーが出る可能性あるかもしれない
func Add(c *gin.Context) {
	var input group.Group
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	userResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if userResult.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: userResult.Err.Error()})
		return
	}

	// check authority
	if userResult.User.Level > 1 {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	if userResult.User.GID != 0 {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "error: You can't create new group", GroupData: nil})
		return
	}

	if err := check(input); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}

	result, err := dbGroup.Create(&group.Group{
		Agree: true, Question: input.Question, Org: input.Org, Status: 0, Bandwidth: input.Bandwidth, Name: input.Name,
		PostCode: input.PostCode, Address: input.Address, Mail: input.Mail, Phone: input.Phone, Country: input.Country,
		Comment: input.Comment, Monitor: input.Monitor, Contract: input.Contract,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error(), GroupData: nil})
		return
	}
	if err := dbUser.Update(user.UpdateStatus, &user.User{Model: gorm.Model{ID: userResult.User.ID}, Status: 10}); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error(), GroupData: nil})
		return
	}
	if err := dbUser.Update(user.UpdateGID, &user.User{Model: gorm.Model{ID: userResult.User.ID}, GID: result.Model.ID}); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, group.Result{Status: true})
	}
}

func Update(c *gin.Context) {
	var input group.Group

	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	authResult := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if authResult.Err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: authResult.Err.Error()})
		return
	}

	if authResult.User.Level > 1 {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "error: failed user level"})
		return
	}
	if authResult.Group.Lock {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "error: This group is locked"})
		return
	}

	data := authResult.Group

	if data.Org != input.Org {
		data.Org = input.Org
	}

	if data.Bandwidth != input.Bandwidth {
		data.Bandwidth = input.Bandwidth
	}
	if data.Name != input.Name {
		data.Name = input.Name
	}
	if data.PostCode != input.PostCode {
		data.PostCode = input.PostCode
	}
	if data.Address != input.Address {
		data.Address = input.Address
	}
	if data.Mail != input.Mail {
		data.Mail = input.Mail
	}
	if data.Phone != input.Phone {
		data.Phone = input.Phone
	}
	if data.Country != input.Country {
		data.Country = input.Country
	}

	if err := dbGroup.Update(group.UpdateInfo, data); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: authResult.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.Result{Status: true})

}

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: result.Err.Error()})
		return
	}

	if result.User.Level >= 10 {
		if result.User.Level > 1 {
			c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "You don't have authority this operation"})
			return
		}
	}

	resultGroup := dbGroup.Get(group.ID, &group.Group{Model: gorm.Model{ID: result.Group.ID}})
	if resultGroup.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: result.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, group.Result{Status: true, GroupData: resultGroup.Group})
}

func GetAll(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, group.ResultAll{Status: false, Error: result.Err.Error()})
		return
	}

	if result.User.Level >= 10 {
		if result.User.Level > 1 {
			c.JSON(http.StatusInternalServerError, group.ResultAll{Status: false, Error: "You don't have authority this operation"})
			return
		}
	}

	resultNetwork := dbNetwork.Get(network.GID, &network.Network{GroupID: result.Group.ID})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, group.ResultAll{Status: false, Error: result.Err.Error()})
		return
	}
	if len(resultNetwork.Network) == 0 {
		resultNetwork.Network = nil
	}

	resultConnection := dbConnection.Get(connection.GID, &connection.Connection{GroupID: result.Group.ID})
	if resultConnection.Err != nil {
		c.JSON(http.StatusInternalServerError, group.ResultAll{Status: false, Error: result.Err.Error()})
		return
	}
	if len(resultConnection.Connection) == 0 {
		resultConnection.Connection = nil
	}

	var resultNetworkUser []networkUser.NetworkUser = nil

	for _, data := range resultNetwork.Network {
		tmp := dbNetworkUser.Get(jpnicUser.GID, &networkUser.NetworkUser{NetworkID: data.ID})
		if tmp.Err != nil {
			c.JSON(http.StatusInternalServerError, group.ResultAll{Status: false, Error: tmp.Err.Error()})
			return
		}
		if len(tmp.NetworkUser) == 0 {
			break
		}
		resultNetworkUser = append(resultNetworkUser, tmp.NetworkUser[0])
	}

	resultJPNICUser := dbJPNICUser.Get(jpnicUser.GID, &jpnicUser.JPNICUser{GroupID: result.Group.ID})
	if resultJPNICUser.Err != nil {
		c.JSON(http.StatusInternalServerError, group.ResultAll{Status: false, Error: result.Err.Error()})
	}
	if len(resultJPNICUser.JPNICUser) == 0 {
		resultJPNICUser.JPNICUser = nil
	}

	c.JSON(http.StatusOK, group.ResultAll{
		Status: true, Group: result.Group, Network: resultNetwork.Network, JPNICUser: resultJPNICUser.JPNICUser,
		NetworkUser: resultNetworkUser, Connection: resultConnection.Connection})
}
