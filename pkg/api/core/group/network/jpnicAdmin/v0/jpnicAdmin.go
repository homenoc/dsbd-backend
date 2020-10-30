package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
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

	log.Println(c.BindJSON(&input))

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, jpnicAdmin.Result{Status: false, Error: result.Err.Error()})
		return
	}

	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusUnauthorized, jpnicAdmin.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	networkResult := dbNetwork.Get(network.GID, &network.Network{GroupID: result.Group.ID})
	if networkResult.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: networkResult.Err.Error()})
		return
	}

	count := 0
	for _, data := range networkResult.Network {
		if data.ID == input.NetworkId {
			count++
		}
	}
	if count == 0 {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: fmt.Sprint("This network id hasn't your group")})
		return
	}

	userResult := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: input.UserId}})
	if userResult.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: networkResult.Err.Error()})
		return
	}

	if userResult.User[0].ID != input.UserId {
		c.JSON(http.StatusBadRequest, jpnicAdmin.Result{Status: false, Error: "This network id hasn't your group"})
		return
	}

	_, err := dbJPNICUser.Create(&jpnicAdmin.JpnicAdmin{NetworkId: input.NetworkId, UserId: input.UserId, Lock: &[]bool{true}[0]})
	if err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: true})
}

func Delete(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, user.Result{Status: false, Error: fmt.Sprintf("id error")})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusBadRequest, jpnicAdmin.Result{Status: false, Error: "wrong id"})
		return
	}

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, jpnicAdmin.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusUnauthorized, jpnicAdmin.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	resultJpnic := dbJPNICUser.Get(jpnicAdmin.ID, &jpnicAdmin.JpnicAdmin{Model: gorm.Model{ID: uint(id)}})
	if resultJpnic.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: resultJpnic.Err.Error()})
		return
	}

	networkResult := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: resultJpnic.Jpnic[0].NetworkId}})
	if networkResult.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: networkResult.Err.Error()})
		return
	}

	if networkResult.Network[0].GroupID != result.Group.ID {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: "group id is not match."})
		return
	}

	if err := dbJPNICUser.Delete(&jpnicAdmin.JpnicAdmin{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: err.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: true})
	}
}

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, jpnicAdmin.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 3 {
		c.JSON(http.StatusUnauthorized, jpnicAdmin.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	networkResult := dbNetwork.Get(network.GID, &network.Network{GroupID: result.Group.ID})
	if networkResult.Err != nil {
		c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: networkResult.Err.Error()})
		return
	}

	var data []jpnicAdmin.JpnicAdmin

	for _, net := range networkResult.Network {
		resultJpnic := dbJPNICUser.Get(jpnicAdmin.NetworkId, &jpnicAdmin.JpnicAdmin{NetworkId: net.ID})
		if resultJpnic.Err != nil {
			c.JSON(http.StatusInternalServerError, jpnicAdmin.Result{Status: false, Error: resultJpnic.Err.Error()})
			return
		}
		for _, detail := range resultJpnic.Jpnic {
			data = append(data, detail)
		}
	}
	c.JSON(http.StatusOK, jpnicAdmin.Result{Status: true, Jpnic: data})
}
