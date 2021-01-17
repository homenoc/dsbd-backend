package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicTech"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbJPNICAdmin "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/jpnicAdmin/v0"
	dbJPNICTech "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/jpnicTech/v0"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input network.Network

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	log.Println(c.BindJSON(&input))

	if _, err := dbNetwork.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, network.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, network.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbNetwork.Delete(&network.Network{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, network.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input network.Network

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, token.Result{Status: false, Error: err.Error()})
		return
	}

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, token.Result{Status: false, Error: err.Error()})
		return
	}

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	result := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusBadRequest, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if err := dbNetwork.Update(network.UpdateAll, replaceAdminNetwork(result.Network[0], input)); err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, network.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, network.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: result.Err.Error()})
		return
	}
	resultJPNICAdmin := dbJPNICAdmin.Get(jpnicAdmin.NetworkId, &jpnicAdmin.JpnicAdmin{NetworkID: uint(id)})
	if resultJPNICAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: resultJPNICAdmin.Err.Error()})
		return
	}
	resultJPNICTech := dbJPNICTech.Get(jpnicTech.NetworkId, &jpnicTech.JpnicTech{NetworkID: uint(id)})
	if resultJPNICTech.Err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: resultJPNICTech.Err.Error()})
		return
	}

	resultUser := dbUser.Get(user.GID, &user.User{GroupID: result.Network[0].GroupID})
	if resultUser.Err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: resultUser.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, network.Result{User: resultUser.User,
		Status: true, Network: result.Network, JPNICAdmin: resultJPNICAdmin.Jpnic, JPNICTech: resultJPNICTech.Jpnic})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbNetwork.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, network.Result{Status: true, Network: result.Network})
	}
}
