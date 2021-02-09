package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	// ID取得
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, network.Result{Status: false, Error: err.Error()})
		return
	}

	// networkIDが0の時エラー処理
	if id == 0 {
		c.JSON(http.StatusBadRequest, network.Result{Status: false, Error: fmt.Sprintf("This id is wrong... ")})
		return
	}

	var input network.Input

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, network.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	log.Println(c.BindJSON(&input))

	// check json
	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, network.Result{Status: false, Error: err.Error()})
		return
	}

	log.Println(input)
	log.Println(input.IP)

	grpIP, err := ipProcess(input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, network.Result{Status: false, Error: err.Error()})
		return
	}

	jh := jpnicHandler{
		admin: input.AdminID, tech: input.TechID, groupID: input.GroupID, jpnicAdmin: nil, jpnicTech: nil,
	}

	// PIアドレスではない場合、jpnic Processを実行
	if !input.PI {
		if err = jh.jpnicProcess(); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, network.Result{Status: false, Error: err.Error()})
			return
		}
	}

	//log.Println(&network.Network{
	//	GroupID: uint(id), Org: input.Org, OrgEn: input.OrgEn, Postcode: input.Postcode, Address: input.Address,
	//	AddressEn: input.AddressEn, RouteV4: input.RouteV4, RouteV6: input.RouteV6, PI: &[]bool{input.PI}[0],
	//	ASN: input.ASN, Open: &[]bool{false}[0], IP: *grpIP, JPNICAdmin: *jh.jpnicAdmin, JPNICTech: *jh.jpnicTech,
	//	V4Name: *input.V4Name, V6Name: *input.V6Name, Lock: &[]bool{input.Lock}[0]})

	//db create for network
	_, err = dbNetwork.Create(&network.Network{
		GroupID: uint(id), Org: input.Org, OrgEn: input.OrgEn, Postcode: input.Postcode, Address: input.Address,
		AddressEn: input.AddressEn, RouteV4: input.RouteV4, RouteV6: input.RouteV6, PI: &[]bool{input.PI}[0],
		ASN: input.ASN, Open: &[]bool{false}[0], IP: *grpIP, JPNICAdmin: *jh.jpnicAdmin, JPNICTech: *jh.jpnicTech,
		V4Name: *input.V4Name, V6Name: *input.V6Name, Lock: &[]bool{input.Lock}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, network.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, network.Result{Status: false, Error: resultAdmin.Err.Error()})
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

	resultUser := dbUser.Get(user.GID, &user.User{GroupID: result.Network[0].GroupID})
	if resultUser.Err != nil {
		c.JSON(http.StatusInternalServerError, network.Result{Status: false, Error: resultUser.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, network.Result{User: resultUser.User, Status: true, Network: result.Network})
}

func Get(c *gin.Context) {
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
