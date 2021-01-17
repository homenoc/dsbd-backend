package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	groupConnection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	groupNetwork "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicTech"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbJPNICAdmin "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/jpnicAdmin/v0"
	dbJPNICTech "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/jpnicTech/v0"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input group.Group

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	log.Println(c.BindJSON(&input))

	if _, err := dbGroup.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, group.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbGroup.Delete(&group.Group{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input group.Group

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	log.Println(c.BindJSON(&input))

	tmp := dbGroup.Get(group.ID, &group.Group{Model: gorm.Model{ID: input.ID}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: tmp.Err.Error()})
		return
	}

	replace, err := updateAdminUser(input, tmp.Group[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, group.Result{Status: false, Error: "error: this email is already registered"})
		return
	}

	if err := dbGroup.Update(group.UpdateAll, replace); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, group.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbGroup.Get(group.ID, &group.Group{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: result.Err.Error()})
		return
	}

	resultUser := dbUser.Get(user.GID, &user.User{GroupID: uint(id)})
	if resultUser.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: result.Err.Error()})
		return
	}

	resultNetwork := dbNetwork.Get(groupNetwork.GID, &groupNetwork.Network{GroupID: uint(id)})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: result.Err.Error()})
		return
	}

	resultConnection := dbConnection.Get(groupConnection.GID, &groupConnection.Connection{GroupID: uint(id)})
	if resultConnection.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: result.Err.Error()})
		return
	}

	var network []group.NetworkInfo

	for _, one := range resultNetwork.Network {
		tmpJPNICAdmin := dbJPNICAdmin.Get(jpnicAdmin.NetworkId, &jpnicAdmin.JpnicAdmin{NetworkID: one.ID})
		if resultConnection.Err != nil {
			c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: result.Err.Error()})
			return
		}
		tmpJPNICTech := dbJPNICTech.Get(jpnicTech.NetworkId, &jpnicTech.JpnicTech{NetworkID: one.ID})
		if resultConnection.Err != nil {
			c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: result.Err.Error()})
			return
		}
		network = append(network, group.NetworkInfo{
			Network:    one,
			JPNICAdmin: tmpJPNICAdmin.Jpnic,
			JPNICTech:  tmpJPNICTech.Jpnic,
		})
	}

	c.JSON(http.StatusOK, group.AdminResult{Status: true, User: resultUser.User, Group: result.Group,
		Network: network, Connection: resultConnection.Connection})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbGroup.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, group.Result{Status: true, Group: result.Group})
	}
}
