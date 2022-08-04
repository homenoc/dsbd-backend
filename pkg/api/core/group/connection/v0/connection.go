package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	connectionTemplate "github.com/homenoc/dsbd-backend/pkg/api/core/template/connection"
	ntt "github.com/homenoc/dsbd-backend/pkg/api/core/template/ntt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	dbConnectionTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/connection/v0"
	dbIPv4RouteTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/ipv4_route/v0"
	dbIPv6RouteTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/ipv6_route/v0"
	dbNTTTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/ntt/v0"
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

	resultConnectionTemplate := dbConnectionTemplate.Get(connectionTemplate.ID,
		&core.ConnectionTemplate{Model: gorm.Model{ID: input.ConnectionTemplateID}})
	if resultConnectionTemplate.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultConnectionTemplate.Err.Error()})
		return
	}

	if *resultConnectionTemplate.Connections[0].NeedInternet {
		resultNTT := dbNTTTemplate.Get(ntt.ID, &core.NTTTemplate{Model: gorm.Model{ID: input.NTTTemplateID}})
		if resultNTT.Err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: resultNTT.Err.Error()})
			return
		}
	}

	// NOCIDが0の時、「どこでも収容」という意味
	if input.NOCID != 0 {
		resultNOC := dbNOC.Get(noc.ID, &core.NOC{Model: gorm.Model{ID: input.NOCID}})
		if resultNOC.Err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: resultNOC.Err.Error()})
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

	// if need_route is true
	if *resultService.Service[0].ServiceTemplate.NeedRoute {
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
			_, err = dbIPv4RouteTemplate.Get(input.IPv4RouteTemplateID)
			if err != nil {
				c.JSON(http.StatusBadRequest, common.Error{Error: "error: invalid ipv4 route template."})
				return
			}
		}
		if ipv6Enable {
			_, err = dbIPv6RouteTemplate.Get(input.IPv6RouteTemplateID)
			if err != nil {
				c.JSON(http.StatusBadRequest, common.Error{Error: "error: invalid ipv4 route template."})
				return
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

	var connectonTemplateID *uint
	if input.ConnectionTemplateID == 0 {
		connectonTemplateID = nil
	} else {
		connectonTemplateID = &[]uint{input.ConnectionTemplateID}[0]
	}

	var IPv4RouteTemplateID *uint
	if input.IPv4RouteTemplateID == 0 {
		IPv4RouteTemplateID = nil
	} else {
		IPv4RouteTemplateID = &[]uint{input.IPv4RouteTemplateID}[0]
	}

	var IPv6RouteTemplateID *uint
	if input.IPv6RouteTemplateID == 0 {
		IPv6RouteTemplateID = nil
	} else {
		IPv6RouteTemplateID = &[]uint{input.IPv6RouteTemplateID}[0]
	}

	var NTTTemplateID *uint
	if input.NTTTemplateID == 0 {
		NTTTemplateID = nil
	} else {
		NTTTemplateID = &[]uint{input.NTTTemplateID}[0]
	}

	var NOCID *uint
	if input.NOCID == 0 {
		NOCID = nil
	} else {
		NOCID = &[]uint{input.NOCID}[0]
	}

	_, err = dbConnection.Create(&core.Connection{
		ServiceID:                resultService.Service[0].ID,
		ConnectionTemplateID:     connectonTemplateID,
		ConnectionComment:        input.ConnectionComment,
		ConnectionNumber:         number,
		IPv4RouteTemplateID:      IPv4RouteTemplateID,
		IPv6RouteTemplateID:      IPv6RouteTemplateID,
		NTTTemplateID:            NTTTemplateID,
		NOCID:                    NOCID,
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
	serviceCode := resultService.Service[0].ServiceTemplate.Type + strconv.Itoa(int(resultService.Service[0].ServiceNumber))
	connectionCodeNew := resultConnectionTemplate.Connections[0].Type + fmt.Sprintf("%03d", number)
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
