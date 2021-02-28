package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"net/http"
)

func GetAdmin(c *gin.Context) {
	adminResult := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if adminResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: adminResult.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, service.Result{Network: config.Conf.Network, Connection: config.Conf.Connection})
}
