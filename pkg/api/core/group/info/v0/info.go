package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/info"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	"net/http"
)

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(0, token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	resultNetwork := dbNetwork.Get(network.Open, &network.Network{GroupID: result.Group.ID})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	var information []info.Info
	var v4 []string
	var v6 []string
	var asn string

	for _, tmpNetwork := range resultNetwork.Network {
		v4 = []string{}
		v6 = []string{}
		asn = tmpNetwork.ASN

		if len(tmpNetwork.IP) > 0 {
			for _, tmpIP := range tmpNetwork.IP {
				if tmpIP.Version == 4 {
					v4 = append(v4, tmpIP.IP)
				} else if tmpIP.Version == 6 {
					v6 = append(v6, tmpIP.IP)
				}
			}
			if len(tmpNetwork.Connection) > 0 {
				for _, tmpConnection := range tmpNetwork.Connection {
					if *tmpConnection.Open {
						information = append(information, info.Info{
							ServiceID: tmpConnection.ServiceID, Service: tmpConnection.Service,
							UserID: tmpConnection.UserID, NOC: tmpConnection.NOC, V4: v4, V6: v6, ASN: asn,
							Assign: tmpConnection.NOCIP, TermIP: tmpConnection.TermIP, NOCIP: tmpConnection.NOCIP,
							LinkV4Our: tmpConnection.LinkV4Our, LinkV4Your: tmpConnection.LinkV4Your,
							LinkV6Our: tmpConnection.LinkV6Our, LinkV6Your: tmpConnection.LinkV6Your,
							Fee: tmpConnection.Fee})
					}
				}
			}
		}
	}

	if len(information) == 0 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "not opening"})
		return
	}

	c.JSON(http.StatusOK, info.Result{Info: information})
}
