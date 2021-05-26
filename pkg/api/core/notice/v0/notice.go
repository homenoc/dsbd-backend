package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	dbNotice "github.com/homenoc/dsbd-backend/pkg/api/store/notice/v0"
	"github.com/jinzhu/gorm"
	"net/http"
)

type noticeHandler struct {
	notice []notice.Notice
}

const layout = "2006-01-02T15:04:05"

//
// DBに入っている情報はUTCベースなので注意が必要
//

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
	var responseNotice []notice.Notice

	h := noticeHandler{notice: responseNotice}

	noticeResult = dbNotice.Get(notice.UIDOrAll, &core.Notice{
		User: []core.User{{Model: gorm.Model{ID: result.User.ID}}},
	})
	if noticeResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	for _, tmpNotice := range noticeResult.Notice {
		h.appendNotice(tmpNotice, tmpNotice.StartTime.Format(layout), tmpNotice.EndTime.Format(layout))
	}

	c.JSON(http.StatusOK, notice.Result{Notice: h.notice})
}

func (h *noticeHandler) appendNotice(data core.Notice, startTime, endTime string) {
	h.notice = append(h.notice, notice.Notice{
		ID:        data.ID,
		Everyone:  *data.Everyone,
		StartTime: startTime,
		EndTime:   endTime,
		Important: *data.Important,
		Fault:     *data.Fault,
		Info:      *data.Info,
		Title:     data.Title,
		Data:      data.Data,
	})
}
