package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnic_user"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbJPNICUser "github.com/homenoc/dsbd-backend/pkg/store/group/network/jpnic_user/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/store/group/v0"
	"github.com/jinzhu/gorm"
	"net/http"
)

func Add(c *gin.Context) {
	var input jpnic_user.JPNICUser
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: result.Err.Error()})
		return
	}

	if err := check(input); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	if err := checkNetworkID(input, result.Group.ID); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: err.Error()})
		return
	}

	_, err := dbJPNICUser.Create(&jpnic_user.JPNICUser{
		GroupID: result.Group.ID, NameJa: input.NameJa, Name: input.Name, OrgJa: input.OrgJa, Org: input.Org,
		PostCode: input.PostCode, AddressJa: input.AddressJa, Address: input.AddressJa, DeptJa: input.DeptJa,
		Dept: input.Dept, PosJa: input.PosJa, Pos: input.Pos, Mail: input.Mail, Tel: input.Tel, Fax: input.Fax,
		OperationID: input.OperationID, TechID: input.TechID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: err.Error()})
		return
	}
	if err := dbGroup.Update(group.UpdateStatus, group.Group{Model: gorm.Model{ID: result.Group.ID},
		Status: result.Group.Status + 1}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, group.Result{Status: true})
}

func Update(c *gin.Context) {
	var input jpnic_user.JPNICUser
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	if err := checkNetworkID(input, result.Group.ID); err != nil {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: err.Error()})
		return
	}

	resultJPNICUser := dbJPNICUser.Get(jpnic_user.ID, &jpnic_user.JPNICUser{Model: gorm.Model{ID: input.ID}})
	if resultJPNICUser.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: result.Err.Error()})
		return
	}
	if resultJPNICUser.JPNICUser[0].GroupID != result.Group.ID {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: "Authentication failure"})
		return
	}

	err := dbJPNICUser.Update(jpnic_user.UpdateInfo, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jpnic_user.Result{Status: false, Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, jpnic_user.Result{Status: true})
	}
}
