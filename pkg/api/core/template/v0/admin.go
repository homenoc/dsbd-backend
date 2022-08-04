package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	template "github.com/homenoc/dsbd-backend/pkg/api/core/template"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbBGPRouter "github.com/homenoc/dsbd-backend/pkg/api/store/noc/bgpRouter/v0"
	dbTunnelEndPointRouterIP "github.com/homenoc/dsbd-backend/pkg/api/store/noc/tunnelEndPointRouterIP/v0"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	dbConnectionTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/connection/v0"
	dbIPv4Template "github.com/homenoc/dsbd-backend/pkg/api/store/template/ipv4/v0"
	dbIPv4RouteTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/ipv4_route/v0"
	dbIPv6Template "github.com/homenoc/dsbd-backend/pkg/api/store/template/ipv6/v0"
	dbIPv6RouteTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/ipv6_route/v0"
	dbMailTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/mail/v0"
	dbNTTTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/ntt/v0"
	dbPaymentCouponTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/payment_coupon/v0"
	dbPaymentMembershipTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/payment_membership/v0"
	dbServiceTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/service/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"net/http"
)

func GetByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
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
	resultIPv4 := dbIPv4Template.GetAll()
	if resultIPv4.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultIPv4.Err.Error()})
		return
	}
	resultIPv6 := dbIPv6Template.GetAll()
	if resultIPv6.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultIPv6.Err.Error()})
		return
	}

	resultIPv4Route, err := dbIPv4RouteTemplate.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	resultIPv6Route, err := dbIPv6RouteTemplate.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
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

	resultPaymentMembership, err := dbPaymentMembershipTemplate.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	resultPaymentCoupon, err := dbPaymentCouponTemplate.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	resultMailTemplate, err := dbMailTemplate.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, template.ResultAdmin{
		Services:                  resultService.Services,
		Connections:               resultConnection.Connections,
		NTTs:                      resultNTT.NTTs,
		NOC:                       resultNOC.NOC,
		BGPRouter:                 resultBGPRouter.BGPRouter,
		TunnelEndPointRouterIP:    resultTunnelEndPointRouterIP.TunnelEndPointRouterIP,
		IPv4:                      resultIPv4.IPv4,
		IPv6:                      resultIPv6.IPv6,
		IPv4Route:                 resultIPv4Route,
		IPv6Route:                 resultIPv6Route,
		User:                      resultUser.User,
		Group:                     resultGroup.Group,
		PaymentMembershipTemplate: resultPaymentMembership,
		PaymentCouponTemplate:     resultPaymentCoupon,
		MailTemplate:              resultMailTemplate,
	})
}
