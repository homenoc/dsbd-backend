package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbJPNICUser "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/jpnicAdmin/v0"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func Add(c *gin.Context) {
	var input jpnicAdmin.JpnicAdmin
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	networkResult := dbNetwork.Get(network.GID, &network.Network{GroupID: result.Group.ID})
	if networkResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: networkResult.Err.Error()})
		return
	}

	count := 0
	for _, data := range networkResult.Network {
		if data.ID == input.NetworkID {
			count++
		}
	}
	if count == 0 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: fmt.Sprint("This network id hasn't your group")})
		return
	}

	userResult := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: input.UserID}})
	if userResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: networkResult.Err.Error()})
		return
	}

	if userResult.User[0].ID != input.UserID {
		c.JSON(http.StatusBadRequest, common.Error{Error: "This network id hasn't your group"})
		return
	}

	_, err = dbJPNICUser.Create(&jpnicAdmin.JpnicAdmin{NetworkID: input.NetworkID, UserID: input.UserID, Lock: &[]bool{true}[0]})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, jpnicAdmin.Result{})
}

func Delete(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("id error")})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: "wrong id"})
		return
	}

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	resultJpnic := dbJPNICUser.Get(jpnicAdmin.ID, &jpnicAdmin.JpnicAdmin{Model: gorm.Model{ID: uint(id)}})
	if resultJpnic.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultJpnic.Err.Error()})
		return
	}

	networkResult := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: resultJpnic.Jpnic[0].NetworkID}})
	if networkResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: networkResult.Err.Error()})
		return
	}

	if networkResult.Network[0].GroupID != result.Group.ID {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "group id is not match."})
		return
	}

	if err := dbJPNICUser.Delete(&jpnicAdmin.JpnicAdmin{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, jpnicAdmin.Result{})
	}
}

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 3 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	networkResult := dbNetwork.Get(network.GID, &network.Network{GroupID: result.Group.ID})
	if networkResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: networkResult.Err.Error()})
		return
	}

	var data []jpnicAdmin.JpnicAdmin

	for _, net := range networkResult.Network {
		resultJpnic := dbJPNICUser.Get(jpnicAdmin.NetworkId, &jpnicAdmin.JpnicAdmin{NetworkID: net.ID})
		if resultJpnic.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: resultJpnic.Err.Error()})
			return
		}
		for _, detail := range resultJpnic.Jpnic {
			data = append(data, detail)
		}
	}
	c.JSON(http.StatusOK, jpnicAdmin.Result{Jpnic: data})
}
