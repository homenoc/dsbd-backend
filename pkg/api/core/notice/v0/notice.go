package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbNotice "github.com/homenoc/dsbd-backend/pkg/store/notice/v0"
	"net/http"
)

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	// Group authentication
	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, notice.Result{Status: false, Error: result.Err.Error()})
		return
	}

	noticeResult := dbNotice.Get(notice.Data, &notice.Notice{UserID: result.User.ID, GroupID: result.Group.ID, Everyone: true})
	if noticeResult.Err != nil {
		c.JSON(http.StatusInternalServerError, notice.Result{Status: false, Error: result.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, notice.Result{Status: true, Notice: noticeResult.Notice})
}
