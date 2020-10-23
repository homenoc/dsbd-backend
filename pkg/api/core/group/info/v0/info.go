package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/info"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	"net/http"
)

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, info.Result{Status: false, Error: result.Err.Error()})
		return
	}

	resultNetwork := dbNetwork.Get(network.GID, &network.Network{GroupID: result.Group.ID})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, info.Result{Status: false, Error: result.Err.Error()})
		return
	}

	resultConnection := dbConnection.Get(connection.GID, &connection.Connection{GroupID: result.Group.ID})
	if resultConnection.Err != nil {
		c.JSON(http.StatusInternalServerError, info.Result{Status: false, Error: result.Err.Error()})
		return
	}

	var information []info.Info

	for _, data := range resultConnection.Connection {
		if *data.Open {
			information = append(information, info.Info{
				ServiceID: data.ServiceID, Service: data.Service, UserID: data.UserId, NOC: data.NOC,
				Assign: data.NOCIP, TermIP: data.TermIP, LinkV4Our: data.LinkV4Our,
				LinkV4Your: data.LinkV4Your, LinkV6Our: data.LinkV6Our, LinkV6Your: data.LinkV6Your, Fee: data.Fee})
		}
	}

	if len(information) == 0 {
		c.JSON(http.StatusInternalServerError, info.Result{Status: false, Error: "not opening"})
		return
	}

	var v4 []string
	var v6 []string
	var asn string

	for i, data := range resultNetwork.Network {
		if data.Open {
			v4 = append(v4, data.V4)
			v6 = append(v6, data.V6)
			if i == 0 {
				asn = data.ASN
			}
			// asnが同じであるかチェック
			if asn != data.ASN {
				c.JSON(http.StatusInternalServerError, info.Result{Status: false, Error: "data mismatch"})
				return
			}
		}
	}

	for i, _ := range information {
		information[i].ASN = asn
		information[i].V4 = v4
		information[i].V6 = v6
	}

	c.JSON(http.StatusOK, info.Result{Status: true, Info: information})
}
