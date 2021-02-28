package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"net/http"
)

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	userResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if userResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: userResult.Err.Error()})
		return
	}

	var network []config.Network
	var connection []config.Connection

	for _, tmp := range config.Conf.Network {
		if !tmp.Hidden {
			network = append(network, tmp)
		}
	}

	for _, tmp := range config.Conf.Connection {
		if !tmp.Hidden {
			connection = append(connection, tmp)
		}
	}

	c.JSON(http.StatusOK, service.Result{Network: network, Connection: connection})
}
