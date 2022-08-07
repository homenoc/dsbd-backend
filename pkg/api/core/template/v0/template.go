package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	template "github.com/homenoc/dsbd-backend/pkg/api/core/template"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	dbConnectionTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/connection/v0"
	dbPaymentCouponTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/payment_coupon/v0"
	dbPaymentMembershipTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/payment_membership/v0"
	dbServiceTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/service/v0"
	"net/http"
)

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	userResult := auth.UserAuthorization(core.Token{UserToken: userToken, AccessToken: accessToken})
	if userResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: userResult.Err.Error()})
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
	resultNOC := dbNOC.GetAll()
	if resultNOC.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNOC.Err.Error()})
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

	c.JSON(http.StatusOK, template.Result{
		Services:                  resultService.Services,
		Connections:               resultConnection.Connections,
		NTTs:                      config.Conf.Template.NTT,
		NOC:                       resultNOC.NOC,
		IPv4:                      config.Conf.Template.V4,
		IPv6:                      config.Conf.Template.V6,
		IPv4Route:                 config.Conf.Template.V4Route,
		IPv6Route:                 config.Conf.Template.V6Route,
		PaymentMembershipTemplate: resultPaymentMembership,
		PaymentCouponTemplate:     resultPaymentCoupon,
	})
}
