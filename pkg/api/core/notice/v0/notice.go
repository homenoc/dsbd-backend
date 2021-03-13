package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	dbNotice "github.com/homenoc/dsbd-backend/pkg/api/store/notice/v0"
	"net/http"
)

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	// Group authentication
	result := auth.UserAuthentication(core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	var noticeResult notice.ResultDatabase

	// Todo #issue #37 Critical
	if result.User.GroupID != 0 {
		noticeResult = dbNotice.Get(notice.GroupID, &core.Notice{
			UserID:   result.User.ID,
			GroupID:  result.User.GroupID,
			Everyone: &[]bool{true}[0],
		})
		if noticeResult.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
			return
		}
	} else {
		noticeResult = dbNotice.Get(notice.UserID, &core.Notice{
			UserID:   result.User.ID,
			Everyone: &[]bool{true}[0],
		})
		if noticeResult.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, notice.Result{Notice: noticeResult.Notice})
}
