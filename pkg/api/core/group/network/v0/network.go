package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func Add(c *gin.Context) {
	var input network.NetworkInput
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	log.Println(c.BindJSON(&input))

	// group authentication
	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, network.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// check user level
	if result.User.Level > 1 {
		c.JSON(http.StatusUnauthorized, network.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	// check json
	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, group.Result{Status: false, Error: err.Error()})
		return
	}

	//log.Println(input)

	var afterStatus uint

	// status check for group
	if result.Group.Status == 2 {
		if input.PI {
			afterStatus = 22
		} else {
			afterStatus = 12
		}
	} else if !(result.Group.Status == 111 || result.Group.Status == 121 || result.Group.Status == 11 || result.Group.Status == 21) {
		c.JSON(http.StatusUnauthorized, network.Result{Status: false, Error: "error: group status"})
		return
	} else {
		// 111,121,11,21の場合はStatusを+1にする
		afterStatus = result.Group.Status + 1
	}

	// db create for network
	net, err := dbNetwork.Create(&network.Network{
		GroupID: result.Group.ID, Org: input.Org, OrgEn: input.OrgEn, Postcode: input.Postcode, Address: input.Address,
		AddressEn: input.AddressEn, Route: input.Route, PI: input.PI, ASN: input.ASN, V4: input.V4, V6: input.V6,
		V4Name: input.V4Name, V6Name: input.V6Name, Date: input.Date, Plan: input.Plan, Lock: false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
		return
	}

	if !input.PI {
		// jpnic Process
		err = jpnicProcess(jpnic{admin: input.AdminID, tech: input.TechID, network: *net})
		if err != nil {
			log.Println(dbNetwork.Delete(&network.Network{Model: gorm.Model{ID: net.ID}}))
		}
	}

	// ---------ここまで処理が通っている場合、DBへの書き込みにすべて成功している
	// GroupのStatusをAfterStatusにする
	if err := dbGroup.Update(group.UpdateStatus, group.Group{Model: gorm.Model{ID: result.Group.ID},
		Status: afterStatus}); err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, network.ResultOne{Status: true, Network: *net})
}

// Todo: 以下の処理は実装中
func Update(c *gin.Context) {
	var input network.Network
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	log.Println(c.BindJSON(&input))

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, network.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusUnauthorized, network.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	if !(result.Group.Status == 211 || result.Group.Status == 221) {
		c.JSON(http.StatusUnauthorized, network.Result{Status: false, Error: "error: group status"})
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
