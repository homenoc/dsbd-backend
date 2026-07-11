package v0

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
	dbNotice "github.com/homenoc/dsbd-backend/pkg/api/store/notice/v0"
)

// GetActive returns the notices currently targeting the caller (GET /notice).
func GetActive(c *gin.Context) {
	u := middleware.CurrentUser(c)

	noticeResult := dbNotice.GetActiveForUser(u.ID)
	if noticeResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: noticeResult.Err.Error()})
		return
	}

	var views []notice.Notice
	for _, n := range noticeResult.Notice {
		views = append(views, notice.Notice{
			StartTime: n.StartTime,
			EndTime:   n.EndTime,
			Everyone:  *n.Everyone,
			Important: *n.Important,
			Fault:     *n.Fault,
			Info:      *n.Info,
			Title:     n.Title,
			Data:      n.Data,
		})
	}
	c.JSON(http.StatusOK, gin.H{"notice": views})
}
