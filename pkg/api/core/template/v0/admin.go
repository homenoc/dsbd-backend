package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	template "github.com/homenoc/dsbd-backend/pkg/api/core/template"
	dbBGPRouter "github.com/homenoc/dsbd-backend/pkg/api/store/noc/bgpRouter/v0"
	dbTunnelEndPointRouterIP "github.com/homenoc/dsbd-backend/pkg/api/store/noc/tunnelEndPointRouterIP/v0"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	dbConnectionTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/connection/v0"
	dbNTTTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/ntt/v0"
	dbServiceTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/service/v0"
	"net/http"
)

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	resultService := dbServiceTemplate.GetAll()
	if resultService.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultService.Err.Error()})
		return
	}
	resultConnection := dbConnectionTemplate.GetAll()
	if resultConnection.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultConnection.Err.Error()})
		return
	}

	resultNTT := dbNTTTemplate.GetAll()
	if resultNTT.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNTT.Err.Error()})
		return
	}

	resultNOC := dbNOC.GetAll()
	if resultNOC.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNOC.Err.Error()})
		return
	}

	resultBGPRouter := dbBGPRouter.GetAll()
	if resultBGPRouter.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultBGPRouter.Err.Error()})
		return
	}

	resultTunnelEndPointRouterIP := dbTunnelEndPointRouterIP.GetAll()
	if resultTunnelEndPointRouterIP.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultTunnelEndPointRouterIP.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, template.Result{
		Services:               resultService.Services,
		Connections:            resultConnection.Connections,
		NTTs:                   resultNTT.NTTs,
		NOC:                    resultNOC.NOC,
		BGPRouter:              resultBGPRouter.BGPRouter,
		TunnelEndPointRouterIP: resultTunnelEndPointRouterIP.TunnelEndPointRouterIP,
	})
}
