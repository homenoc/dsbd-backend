package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/jpnicTech"
	dbJPNICTech "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/jpnicTech/v0"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func Add(c *gin.Context) {
	var input core.JPNICTech
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// check authority
	if result.User.Level > 2 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	networkResult := dbService.Get(service.GID, &core.Service{GroupID: result.User.GroupID})
	if networkResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: networkResult.Err.Error()})
		return
	}

	_, err = dbJPNICTech.Create(&core.JPNICTech{
		Org:       input.Org,
		OrgEn:     input.OrgEn,
		PostCode:  input.PostCode,
		Address:   input.Address,
		AddressEn: input.AddressEn,
		Dept:      input.Dept,
		DeptEn:    input.DeptEn,
		Tel:       input.Tel,
		Fax:       input.Fax,
		Country:   input.Country,
		Lock:      &[]bool{true}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, jpnicTech.Result{})
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

	result := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 2 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	resultTech := dbJPNICTech.Get(jpnicTech.ID, &core.JPNICTech{Model: gorm.Model{ID: uint(id)}})
	if resultTech.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultTech.Err.Error()})
		return
	}

	if err = dbJPNICTech.Delete(&core.JPNICTech{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, jpnicTech.Result{})
	}
}
