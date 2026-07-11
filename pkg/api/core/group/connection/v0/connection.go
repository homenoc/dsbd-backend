package v0

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
)

func Add(c *gin.Context) {
	var input connection.Input

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

	user := middleware.CurrentUser(c)

	// check authority
	if !core.CanManageServices(user.Level) {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	// status check for group
	if !*user.Group.Pass {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: Your group has not yet been reviewed."})
		return
	}

	if *user.Group.ExpiredStatus != core.ExpiredNone {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: failed group status"})
		return
	}

	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// check input.ConnectionType and getting connection template
	connectionTemplate, err := core.GetConnectionType(input.ConnectionType)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// check preferredAP and NTT (internet)
	if connectionTemplate.NeedInternet {
		err = config.CheckIncludePreferredAPTemplate(input.PreferredAP)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}

		err = config.CheckIncludeNTTTemplate(input.NTT)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	// check IX fields (IXP connection)
	if input.ConnectionType == "IXP" {
		if err = config.CheckIXFields(input.IX, input.IXPeerType, input.IXVlanID); err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	resultService := dbService.GetByID(uint(id))
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
	if resultService.Service[0].GroupID != user.Group.ID {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: GroupID does not match."})
		return
	}

	// getting service with template
	resultServiceWithTemplate, err := core.GetServiceType(resultService.Service[0].ServiceType)
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

	resultConnection := dbConnection.GetByServiceID(uint(id))
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
		IX:                       input.IX,
		IXPeerType:               input.IXPeerType,
		IXVlanID:                 input.IXVlanID,
		LinkV4Your:               input.LinkV4Your,
		LinkV6Your:               input.LinkV6Your,
		IPv4Route:                input.IPv4Route,
		IPv6Route:                input.IPv6Route,
		NTT:                      input.NTT,
		PreferredAP:              input.PreferredAP,
		BGPRouterID:              nil,
		TunnelEndPointRouterIPID: nil,
		TermIP:                   input.TermIP,
		RFC8950:                  input.RFC8950,
		Address:                  input.Address,
		Comment:                  input.Comment,
		Monitor:                  &[]bool{input.Monitor}[0],
		Enable:                   &[]bool{true}[0],
		Open:                     &[]bool{false}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	applicant := "[" + strconv.Itoa(int(user.ID)) + "] " + user.Name + "(" + user.NameEn + ")"
	groupName := "[" + strconv.Itoa(int(user.Group.ID)) + "] " + user.Group.Org + "(" + user.Group.OrgEn + ")"
	serviceCode := resultServiceWithTemplate.Type + strconv.Itoa(int(resultService.Service[0].ServiceNumber))
	connectionCodeNew := connectionTemplate.Type + fmt.Sprintf("%03d", number)
	connectionCodeComment := input.ConnectionComment

	noticeAdd(applicant, groupName, serviceCode, connectionCodeNew, connectionCodeComment)

	//if err = dbGroup.Update(group.UpdateStatus, core.Group{
	//	Model:  gorm.Model{ID: user.Group.ID},
	//	Status: &[]uint{4}[0],
	//}); err != nil {
	//	c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	//	return
	//}

	if err = dbService.UpdateAddAllow(resultService.Service[0].ID, false); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	notification.NoticeUpdateStatus(groupName, core.StatusOpening.Label(),
		core.TransitionText(core.StatusConnectionInput, core.StatusOpening))

	c.JSON(http.StatusOK, common.Result{})
}
