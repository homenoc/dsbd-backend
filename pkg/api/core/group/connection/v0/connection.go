package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func Add(c *gin.Context) {
	var input connection.Input
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// ID取得
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: ID is 0"})
		return
	}

	result := auth.GroupAuthorization(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 2 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	// status check for group
	if !*result.User.Group.Pass {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: Your group has not yet been reviewed."})
		return
	}

	if *result.User.Group.ExpiredStatus != 0 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: failed group status"})
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

	// check preferredAP
	err = config.CheckIncludePreferredAPTemplate(input.PreferredAP)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// check NTT (internet)
	if connectionTemplate.NeedInternet {
		err = config.CheckIncludeNTTTemplate(input.NTT)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	resultService := dbService.Get(service.ID, &core.Service{Model: gorm.Model{ID: uint(id)}})
	if resultService.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultService.Err.Error()})
		return
	}

	// check service enable
	if !*resultService.Service[0].Enable {
		c.JSON(http.StatusBadRequest, common.Error{Error: "You don't allow this operation. [enable]"})
		return
	}

	// check service pass
	if !*resultService.Service[0].Pass {
		c.JSON(http.StatusBadRequest, common.Error{Error: "You don't allow this operation. [pass]"})
		return
	}

	// check add_allow
	if !*resultService.Service[0].AddAllow {
		c.JSON(http.StatusBadRequest, common.Error{Error: "You don't allow this operation. [add_allow]"})
		return
	}

	// GroupIDが一致しない場合はエラーを返す
	if resultService.Service[0].GroupID != result.User.Group.ID {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: GroupID does not match."})
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
		ipv4Enable := false
		ipv6Enable := false

		for _, tmpServiceIP := range resultService.Service[0].IP {
			if tmpServiceIP.Version == 4 {
				ipv4Enable = true
				break
			}
			if tmpServiceIP.Version == 6 {
				ipv6Enable = true
				break
			}
		}

		if ipv4Enable {
			err = config.CheckIncludeV4RouteTemplate(input.IPv4Route)
			if err != nil {
				c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			}
		}

		if ipv6Enable {
			err = config.CheckIncludeV6RouteTemplate(input.IPv6Route)
			if err != nil {
				c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			}
		}
	}

	resultConnection := dbConnection.Get(connection.ServiceID, &core.Connection{ServiceID: uint(id)})
	if resultConnection.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultConnection.Err.Error()})
		return
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
		ServiceID:                resultService.Service[0].ID,
		ConnectionType:           input.ConnectionType,
		ConnectionComment:        input.ConnectionComment,
		ConnectionNumber:         number,
		IPv4Route:                input.IPv4Route,
		IPv6Route:                input.IPv6Route,
		NTT:                      input.NTT,
		PreferredAP:              input.PreferredAP,
		BGPRouterID:              nil,
		TunnelEndPointRouterIPID: nil,
		TermIP:                   input.TermIP,
		Address:                  input.Address,
		Monitor:                  &[]bool{input.Monitor}[0],
		Enable:                   &[]bool{true}[0],
		Open:                     &[]bool{false}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	applicant := "[" + strconv.Itoa(int(result.User.ID)) + "] " + result.User.Name + "(" + result.User.NameEn + ")"
	groupName := "[" + strconv.Itoa(int(result.User.Group.ID)) + "] " + result.User.Group.Org + "(" + result.User.Group.OrgEn + ")"
	serviceCode := resultServiceWithTemplate.Type + strconv.Itoa(int(resultService.Service[0].ServiceNumber))
	connectionCodeNew := connectionTemplate.Type + fmt.Sprintf("%03d", number)
	connectionCodeComment := input.ConnectionComment

	noticeAdd(applicant, groupName, serviceCode, connectionCodeNew, connectionCodeComment)

	//if err = dbGroup.Update(group.UpdateStatus, core.Group{
	//	Model:  gorm.Model{ID: result.User.Group.ID},
	//	Status: &[]uint{4}[0],
	//}); err != nil {
	//	c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	//	return
	//}

	if err = dbService.Update(service.UpdateAll, core.Service{
		Model:    gorm.Model{ID: resultService.Service[0].ID},
		AddAllow: &[]bool{false}[0],
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	notification.NoticeUpdateStatus(groupName, "開通作業中", "3[接続情報記入段階(User)] =>4[開通作業中]")

	c.JSON(http.StatusOK, common.Result{})
}
