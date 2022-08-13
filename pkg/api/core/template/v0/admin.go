package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	template "github.com/homenoc/dsbd-backend/pkg/api/core/template"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbBGPRouter "github.com/homenoc/dsbd-backend/pkg/api/store/noc/bgpRouter/v0"
	dbTunnelEndPointRouterIP "github.com/homenoc/dsbd-backend/pkg/api/store/noc/tunnelEndPointRouterIP/v0"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"net/http"
)

func GetByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
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

	resultUser := dbUser.GetAll()
	if resultUser.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultUser.Err.Error()})
		return
	}

	resultGroup := dbGroup.GetAll()
	if resultGroup.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultGroup.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, template.ResultAdmin{
		Services:               config.Conf.Template.Service,
		Connections:            config.Conf.Template.Connection,
		NTTs:                   config.Conf.Template.NTT,
		NOC:                    resultNOC.NOC,
		BGPRouter:              resultBGPRouter.BGPRouter,
		TunnelEndPointRouterIP: resultTunnelEndPointRouterIP.TunnelEndPointRouterIP,
		IPv4:                   config.Conf.Template.V4,
		IPv6:                   config.Conf.Template.V6,
		IPv4Route:              config.Conf.Template.V4Route,
		IPv6Route:              config.Conf.Template.V6Route,
		User:                   resultUser.User,
		Group:                  resultGroup.Group,
		PaymentMembership:      config.Conf.Template.Membership,
		MailTemplate:           config.Conf.Template.Mail,
		PreferredAP:            config.Conf.Template.PreferredAP,
		MemberType: []core.ConstantMembership{
			core.MemberTypeStandard,
			core.MemberTypeStudent,
			core.MemberTypeCommittee,
			core.MemberTypeDisable,
		},
	})
}
