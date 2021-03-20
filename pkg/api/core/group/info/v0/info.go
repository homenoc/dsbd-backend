package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/info"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	"net/http"
	"strconv"
)

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	if !(0 < result.User.Level && result.User.Level <= 3) {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: access is not permitted"})
		return
	}

	resultService := dbService.Get(service.Open, &core.Service{GroupID: result.User.GroupID})
	if resultService.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultService.Err.Error()})
		return
	}

	if len(resultService.Service) == 0 {
		c.JSON(http.StatusNoContent, common.Error{Error: "error: No service information available."})
		return
	}

	var infoInterface []info.Info

	for _, tmpService := range resultService.Service {
		if *tmpService.Open {
			for _, tmpConnection := range tmpService.Connection {
				var fee string
				var v4, v6 []string
				if *tmpService.Fee == 0 {
					fee = "Free"
				} else {
					fee = strconv.Itoa(int(*tmpService.Fee)) + "円"
				}
				serviceID := strconv.Itoa(int(tmpService.GroupID)) + "-" + tmpService.ServiceTemplate.Type +
					fmt.Sprintf("%03d", tmpService.ServiceNumber) + "-" + tmpConnection.ConnectionTemplate.Type +
					fmt.Sprintf("%03d", tmpConnection.ConnectionNumber)

				for _, tmpIP := range tmpService.IP {
					if tmpIP.Version == 4 {
						v4 = append(v4, tmpIP.IP)
					} else if tmpIP.Version == 6 {
						v6 = append(v6, tmpIP.IP)
					}
				}

				if *tmpConnection.Open {
					infoInterface = append(infoInterface, info.Info{
						ServiceID:  serviceID,
						Service:    tmpService.ServiceTemplate.Name,
						Assign:     *tmpService.ServiceTemplate.NeedJPNIC,
						ASN:        *tmpService.ASN,
						V4:         v4,
						V6:         v6,
						NOC:        tmpConnection.NOC.Name,
						NOCIP:      tmpConnection.TunnelEndPointRouterIP.IP,
						TermIP:     tmpConnection.TermIP,
						LinkV4Our:  tmpConnection.LinkV4Our,
						LinkV4Your: tmpConnection.LinkV4Your,
						LinkV6Our:  tmpConnection.LinkV6Our,
						LinkV6Your: tmpConnection.LinkV6Your,
						Fee:        fee,
					})
				}
			}
		}
	}

	c.JSON(http.StatusOK, info.Result{Info: infoInterface})
}
