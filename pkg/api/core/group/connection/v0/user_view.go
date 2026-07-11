package v0

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
)

// loadUserGroup resolves the caller's group graph for the user-facing reads.
// ok=false means a response has been written. hasGroup=false with ok=true
// means the caller has no group yet (respond with an empty collection).
func loadUserGroup(c *gin.Context) (g core.Group, hasGroup bool, ok bool) {
	u := middleware.CurrentUser(c)
	if u.GroupID == nil {
		return core.Group{}, false, true
	}
	if !core.CanViewGroup(u.Level) {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: access is not permitted"})
		return core.Group{}, false, false
	}
	groupResult := dbGroup.GetByID(*u.GroupID)
	if groupResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: groupResult.Err.Error()})
		return core.Group{}, false, false
	}
	if len(groupResult.Group) == 0 {
		return core.Group{}, false, true
	}
	return groupResult.Group[0], true, true
}

// Get returns the group's enabled connections as seen by users (GET /connection).
func Get(c *gin.Context) {
	g, hasGroup, ok := loadUserGroup(c)
	if !ok {
		return
	}
	if !hasGroup {
		c.JSON(http.StatusOK, gin.H{"connection": []connection.Connection{}})
		return
	}

	var views []connection.Connection
	for _, s := range g.Services {
		if _, err := core.GetServiceType(s.ServiceType); err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}
		serviceCode := core.ServiceCode(s.GroupID, s.ServiceType, s.ServiceNumber)
		for _, conn := range s.Connection {
			if !*conn.Enable {
				continue
			}
			views = append(views, connection.Connection{
				ID:        conn.ID,
				ServiceID: core.ConnectionCode(serviceCode, conn.ConnectionType, conn.ConnectionNumber),
				Open:      *conn.Open,
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{"connection": views})
}

// GetInfo returns the derived network summary of the group's opened
// connections (GET /info) — the contract-disclosure data the web Info page
// renders.
func GetInfo(c *gin.Context) {
	g, hasGroup, ok := loadUserGroup(c)
	if !ok {
		return
	}
	if !hasGroup {
		c.JSON(http.StatusOK, gin.H{"info": []connection.Info{}})
		return
	}

	var infos []connection.Info
	for _, s := range g.Services {
		serviceType, err := core.GetServiceType(s.ServiceType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}
		if !(*s.Pass && *s.Enable) {
			continue
		}

		var v4, v6 []string
		for _, ip := range s.IP {
			if !*ip.Open {
				continue
			}
			switch ip.Version {
			case 4:
				v4 = append(v4, ip.IP)
			case 6:
				v6 = append(v6, ip.IP)
			}
		}

		serviceCode := core.ServiceCode(s.GroupID, s.ServiceType, s.ServiceNumber)
		for _, conn := range s.Connection {
			if !(*conn.Open && *conn.Enable) {
				continue
			}
			var asn uint = 0
			if s.ASN != nil {
				asn = *s.ASN
			}
			infos = append(infos, connection.Info{
				ServiceID:  core.ConnectionCode(serviceCode, conn.ConnectionType, conn.ConnectionNumber),
				Service:    serviceType.Name,
				Assign:     serviceType.NeedJPNIC,
				ASN:        asn,
				V4:         v4,
				V6:         v6,
				Fee:        "Free",
				NOC:        conn.BGPRouter.NOC.Name,
				NOCIP:      conn.TunnelEndPointRouterIP.IP,
				TermIP:     conn.TermIP,
				RFC8950:    conn.RFC8950,
				LinkV4Our:  conn.LinkV4Our,
				LinkV4Your: conn.LinkV4Your,
				LinkV6Our:  conn.LinkV6Our,
				LinkV6Your: conn.LinkV6Your,
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{"info": infos})
}
