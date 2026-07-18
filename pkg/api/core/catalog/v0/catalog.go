package v0

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/catalog"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
)

// build assembles the catalog. hideHidden drops services flagged Hidden (user
// audience); admin sees everything.
func build(hideHidden bool) catalog.Result {
	services := core.ServiceTypes()
	if hideHidden {
		filtered := make([]core.ServiceType, 0, len(services))
		for _, s := range services {
			if !s.Hidden {
				filtered = append(filtered, s)
			}
		}
		services = filtered
	}
	return catalog.Result{
		Services:          services,
		Connections:       core.ConnectionTypes(),
		IX:                config.Conf.Template.IX,
		NTTs:              config.Conf.Template.NTT,
		IPv4:              config.Conf.Template.V4,
		IPv6:              config.Conf.Template.V6,
		IPv4Route:         config.Conf.Template.V4Route,
		IPv6Route:         config.Conf.Template.V6Route,
		PreferredAP:       config.Conf.Template.PreferredAP,
		PaymentMembership: config.Conf.Template.Membership,
		MailTemplate:      config.Conf.Template.Mail,
		MemberType:        core.MemberTypes,
	}
}

// Get serves the catalog to authenticated users (Hidden service types omitted).
func Get(c *gin.Context) {
	_ = middleware.CurrentUser(c)
	c.JSON(http.StatusOK, build(true))
}

// GetByAdmin serves the full catalog to the admin app.
func GetByAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, build(false))
}
