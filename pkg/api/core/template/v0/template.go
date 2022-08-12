package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	template "github.com/homenoc/dsbd-backend/pkg/api/core/template"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
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

	var resultService []config.ServiceTemplate
	for _, serviceTemplate := range config.Conf.Template.Service {
		if !serviceTemplate.Hidden {
			resultService = append(resultService, serviceTemplate)
		}
	}
	resultNOC := dbNOC.GetAll()
	if resultNOC.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNOC.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, template.Result{
		Services:                  resultService,
		Connections:               config.Conf.Template.Connection,
		NTTs:                      config.Conf.Template.NTT,
		NOC:                       resultNOC.NOC,
		IPv4:                      config.Conf.Template.V4,
		IPv6:                      config.Conf.Template.V6,
		IPv4Route:                 config.Conf.Template.V4Route,
		IPv6Route:                 config.Conf.Template.V6Route,
		PaymentMembershipTemplate: config.Conf.Template.Membership,
	})
}
