package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	jpnicAdmin "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	jpnicTech "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicTech"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/store/group/connection/v0"
	dbJpnicAdmin "github.com/homenoc/dsbd-backend/pkg/store/group/network/jpnicAdmin/v0"
	dbJpnicTech "github.com/homenoc/dsbd-backend/pkg/store/group/network/jpnicTech/v0"
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
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "error: You can't create new group", Group: nil})
		return
	}

	if err := check(input); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}

	result, err := dbGroup.Create(&group.Group{
		Agree: true, Question: input.Question, Org: input.Org, Status: 0, Bandwidth: input.Bandwidth,
		Comment: input.Comment, Contract: input.Contract,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error(), Group: nil})
		return
	}
	if err := dbUser.Update(user.UpdateGID, &user.User{Model: gorm.Model{ID: userResult.User.ID}, GID: result.Model.ID}); err != nil {
		dbGroup.Delete(&group.Group{Model: gorm.Model{ID: result.ID}})
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

	c.JSON(http.StatusOK, group.Result{Status: true, Group: resultGroup.Group})
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

	var resultJpnicTech []jpnicTech.JpnicTech = nil
	var resultJpnicAdmin []jpnicAdmin.JpnicAdmin = nil

	for _, data := range resultNetwork.Network {
		tmpAdmin := dbJpnicAdmin.Get(jpnicAdmin.NetworkId, &jpnicAdmin.JpnicAdmin{NetworkId: data.ID})
		if tmpAdmin.Err != nil {
			c.JSON(http.StatusInternalServerError, group.ResultAll{Status: false, Error: tmpAdmin.Err.Error()})
			return
		}
		if len(tmpAdmin.Jpnic) == 0 {
			break
		}
		resultJpnicAdmin = append(resultJpnicAdmin, tmpAdmin.Jpnic[0])

		tmpTech := dbJpnicTech.Get(jpnicAdmin.NetworkId, &jpnicTech.JpnicTech{NetworkId: data.ID})
		if tmpAdmin.Err != nil {
			c.JSON(http.StatusInternalServerError, group.ResultAll{Status: false, Error: tmpAdmin.Err.Error()})
			return
		}
		if len(tmpTech.Jpnic) == 0 {
			break
		}
		for _, tmpTechDetail := range tmpTech.Jpnic {
			resultJpnicTech = append(resultJpnicTech, tmpTechDetail)
		}
	}

	c.JSON(http.StatusOK, group.ResultAll{
		Status: true, Group: result.Group, Network: resultNetwork.Network, JpnicAdmin: resultJpnicAdmin,
		JpnicTech: resultJpnicTech, Connection: resultConnection.Connection})
}
