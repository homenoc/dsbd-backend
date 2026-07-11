package v0

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
)

// Get returns the group's enabled services as seen by users (GET /service).
func Get(c *gin.Context) {
	u := middleware.CurrentUser(c)
	if u.GroupID == nil {
		c.JSON(http.StatusOK, gin.H{"service": []service.UserView{}})
		return
	}
	if !core.CanViewGroup(u.Level) {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: access is not permitted"})
		return
	}
	groupResult := dbGroup.GetByID(*u.GroupID)
	if groupResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: groupResult.Err.Error()})
		return
	}
	if len(groupResult.Group) == 0 {
		c.JSON(http.StatusOK, gin.H{"service": []service.UserView{}})
		return
	}

	views, err := userViews(groupResult.Group[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"service": views})
}

// userViews projects a group's enabled services onto their user-facing shape.
func userViews(g core.Group) ([]service.UserView, error) {
	var out []service.UserView
	for _, s := range g.Services {
		serviceType, err := core.GetServiceType(s.ServiceType)
		if err != nil {
			return nil, err
		}
		if !*s.Enable {
			continue
		}

		// The admin contact intentionally omits its address fields on the wire.
		admin := service.UserJPNICView{
			ID:       s.JPNICAdmin.ID,
			Org:      s.JPNICAdmin.Org,
			OrgEn:    s.JPNICAdmin.OrgEn,
			PostCode: s.JPNICAdmin.PostCode,
			Name:     s.JPNICAdmin.Name,
			NameEn:   s.JPNICAdmin.NameEn,
			Dept:     s.JPNICAdmin.Dept,
			DeptEn:   s.JPNICAdmin.DeptEn,
			Tel:      s.JPNICAdmin.Tel,
			Fax:      s.JPNICAdmin.Fax,
			Mail:     s.JPNICAdmin.Mail,
			Country:  s.JPNICAdmin.Country,
		}

		var tech []service.UserJPNICView
		for _, t := range s.JPNICTech {
			tech = append(tech, service.UserJPNICView{
				ID:        t.ID,
				Name:      t.Name,
				NameEn:    t.NameEn,
				Org:       t.Org,
				OrgEn:     t.OrgEn,
				PostCode:  t.PostCode,
				Address:   t.Address,
				AddressEn: t.AddressEn,
				Dept:      t.Dept,
				DeptEn:    t.DeptEn,
				Tel:       t.Tel,
				Fax:       t.Fax,
				Mail:      t.Mail,
				Country:   t.Country,
			})
		}

		var ips []service.UserIPView
		for _, ip := range s.IP {
			if !*ip.Open {
				continue
			}
			var plans []service.UserPlanView
			for _, plan := range ip.Plan {
				plans = append(plans, service.UserPlanView{
					ID:       plan.ID,
					IPID:     plan.IPID,
					Name:     plan.Name,
					After:    plan.After,
					HalfYear: plan.HalfYear,
					OneYear:  plan.OneYear,
				})
			}
			ips = append(ips, service.UserIPView{
				ID:      ip.ID,
				Version: ip.Version,
				Name:    ip.Name,
				IP:      ip.IP,
				Plan:    plans,
				UseCase: ip.UseCase,
			})
		}

		out = append(out, service.UserView{
			ID:             s.ID,
			ServiceID:      core.ServiceCode(s.GroupID, s.ServiceType, s.ServiceNumber),
			ServiceType:    s.ServiceType,
			NeedRoute:      serviceType.NeedRoute,
			NeedBGP:        serviceType.NeedBGP,
			NeedJPNIC:      serviceType.NeedJPNIC,
			AddAllow:       *s.AddAllow,
			Pass:           *s.Pass,
			Org:            s.Org,
			OrgEn:          s.OrgEn,
			PostCode:       s.PostCode,
			Address:        s.Address,
			AddressEn:      s.AddressEn,
			ASN:            s.ASN,
			AveUpstream:    s.AveUpstream,
			MaxUpstream:    s.MaxUpstream,
			AveDownstream:  s.AveDownstream,
			MaxDownstream:  s.MaxDownstream,
			MaxBandWidthAS: s.MaxBandWidthAS,
			IP:             ips,
			JPNICAdmin:     admin,
			JPNICTech:      tech,
		})
	}
	return out, nil
}
