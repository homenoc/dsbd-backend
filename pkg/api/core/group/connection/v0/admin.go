package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddByAdmin(c *gin.Context) {
	// ID取得
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// IDが0の時エラー処理
	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("This id is wrong... ")})
		return
	}

	var input connection.Input

	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// check input.ConnectionType and getting connection template
	connectionTemplate, err := config.GetConnectionTemplate(input.ConnectionType)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// check NTT(Internet)
	err = config.CheckIncludeNTTTemplate(input.NTT)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// check preferredAP
	err = config.CheckIncludePreferredAPTemplate(input.PreferredAP)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultNOC := dbNOC.Get(noc.ID, &core.NOC{Model: gorm.Model{ID: input.NOCID}})
	if resultNOC.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultNOC.Err.Error()})
		return
	}

	resultService := dbService.Get(connection.ID, &core.Service{Model: gorm.Model{ID: uint(id)}})
	if resultService.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultService.Err.Error()})
		return
	}

	resultConnection := dbConnection.Get(connection.ServiceID, &core.Connection{ServiceID: uint(id)})
	if resultConnection.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultConnection.Err.Error()})
		return
	}

	// getting service with template
	resultServiceWithTemplate, err := config.GetServiceTemplate(resultService.Service[0].ServiceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	// if need_route is true
	if resultServiceWithTemplate.NeedRoute {
		err = config.CheckIncludeV4RouteTemplate(input.IPv4Route)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}

		err = config.CheckIncludeV6RouteTemplate(input.IPv6Route)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	var number uint = 1
	for _, tmp := range resultConnection.Connection {
		if tmp.ConnectionNumber >= 1 {
			number = tmp.ConnectionNumber + 1
		}
	}

	if number >= 999 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "error: over number"})
		return
	}

	_, err = dbConnection.Create(&core.Connection{
		ServiceID:         resultService.Service[0].ID,
		ConnectionType:    input.ConnectionType,
		ConnectionComment: input.ConnectionComment,
		ConnectionNumber:  number,
		IPv4Route:         input.IPv4Route,
		IPv6Route:         input.IPv6Route,
		NTT:               input.NTT,
		PreferredAP:       input.PreferredAP,
		NOCID:             &[]uint{input.NOCID}[0],
		TermIP:            input.TermIP,
		Address:           input.Address,
		Monitor:           &[]bool{input.Monitor}[0],
		Enable:            &[]bool{true}[0],
		Open:              &[]bool{false}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	serviceCode := resultServiceWithTemplate.Type + strconv.Itoa(int(resultService.Service[0].ServiceNumber))
	connectionCodeNew := connectionTemplate.Type + fmt.Sprintf("%03d", number)
	connectionCodeComment := input.ConnectionComment
	noticeAdd("", strconv.Itoa(id), serviceCode, connectionCodeNew, connectionCodeComment)

	c.JSON(http.StatusOK, connection.Result{})
}

func DeleteByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbConnection.Delete(&core.Connection{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, connection.Result{})
}

func UpdateByAdmin(c *gin.Context) {
	var input core.Connection

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	tmp := dbConnection.Get(connection.ID, &core.Connection{Model: gorm.Model{ID: uint(id)}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: tmp.Err.Error()})
		return
	}

	noticeUpdateByAdmin(tmp.Connection[0], input)

	input.ID = uint(id)

	if err = dbConnection.Update(connection.UpdateAll, input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, connection.Result{})
}

func GetByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := dbConnection.Get(connection.ID, &core.Connection{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, connection.Result{Connection: result.Connection})
}

func GetAllByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbConnection.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, connection.Result{Connection: result.Connection})
	}
}
