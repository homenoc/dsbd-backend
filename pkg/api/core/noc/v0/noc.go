package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	"net/http"
)

func GetAll(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	userResult := auth.UserAuthentication(core.Token{UserToken: userToken, AccessToken: accessToken})
	if userResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: userResult.Err.Error()})
		return
	}

	result := dbNOC.GetAll()
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	var nocTmp noc.ResultAllUser

	for _, tmp := range result.NOC {
		if *tmp.Enable {
			nocTmp.NOC = append(nocTmp.NOC, noc.ResultOneUser{
				ID:       tmp.ID,
				Name:     tmp.Name,
				Location: tmp.Location,
			})
		}
	}

	c.JSON(http.StatusOK, nocTmp)
}
