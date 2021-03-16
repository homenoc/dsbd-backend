package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	dbNotice "github.com/homenoc/dsbd-backend/pkg/api/store/notice/v0"
	"net/http"
	"strconv"
)

type noticeHandler struct {
	notice []notice.Notice
}

const layout = "2006-01-02 15:04:05"

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

	if result.User.GroupID != 0 {

		serviceResult := dbService.Get(service.Open, &core.Service{GroupID: result.User.GroupID})
		if serviceResult.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
			return
		}

		var nocIDs []string

		for _, tmpService := range serviceResult.Service {
			for _, tmpConnection := range tmpService.Connection {
				if tmpConnection.BGPRouter.NOCID != 0 {
					if !arrayContains(nocIDs, strconv.Itoa(int(tmpConnection.BGPRouter.NOCID))) {
						nocIDs = append(nocIDs, strconv.Itoa(int(tmpConnection.BGPRouter.NOCID)))
					}
				}
			}
		}

		noticeResult = dbNotice.GetArray(notice.UIDOrGIDOrNOCAllOrAll, &core.Notice{
			UserID:  result.User.ID,
			GroupID: result.User.GroupID,
		}, nocIDs)
		if noticeResult.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
			return
		}

		for _, tmpNotice := range noticeResult.Notice {
			h.appendNotice(tmpNotice, tmpNotice.StartTime.Format(layout), tmpNotice.EndTime.Format(layout))
		}

	} else {
		noticeResult = dbNotice.Get(notice.UIDOrAll, &core.Notice{
			UserID: result.User.ID,
		})
		if noticeResult.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
			return
		}
		for _, tmpNotice := range noticeResult.Notice {
			h.appendNotice(tmpNotice, tmpNotice.StartTime.Format(layout), tmpNotice.EndTime.Format(layout))
		}
	}

	c.JSON(http.StatusOK, notice.Result{Notice: h.notice})
}

func arrayContains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func (h *noticeHandler) appendNotice(data core.Notice, startTime, endTime string) {
	h.notice = append(h.notice, notice.Notice{
		ID:        data.ID,
		UserID:    data.UserID,
		GroupID:   data.GroupID,
		NOCID:     data.NOCID,
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
