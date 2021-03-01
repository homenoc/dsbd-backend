package v0

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func Add(c *gin.Context) {
	var input connection.Connection
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := auth.GroupAuthentication(0, token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 1 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	// status check for group
	if !(*result.Group.Status == 3 && *result.Group.ExpiredStatus == 0 && *result.Group.Pass) {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: failed group status"})
		return
	}

	if err = check(result.Group.ID, true, input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultConnection := dbConnection.Get(connection.GID, &connection.Connection{GroupID: result.Group.ID})
	if resultConnection.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultConnection.Err.Error()})
		return
	}
	var number uint = 1
	for _, tmp := range resultConnection.Connection {
		if tmp.ConnectionNumber >= 1 {
			number = tmp.ConnectionNumber + 1
		}
	}

	if number >= 999 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "error: over number"})
		return
	}

	_, err = dbConnection.Create(&connection.Connection{
		GroupID:          result.Group.ID,
		UserID:           input.UserID,
		ConnectionType:   input.ConnectionType,
		ConnectionNumber: number,
		NTT:              input.NTT,
		NOC:              input.NOC,
		TermIP:           input.TermIP,
		Monitor:          input.Monitor,
		Open:             &[]bool{false}[0],
		Lock:             &[]bool{true}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if err = dbGroup.Update(group.UpdateStatus, group.Group{
		Model:  gorm.Model{ID: result.Group.ID},
		Status: &[]uint{4}[0],
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "接続情報登録"}).
		AddField(slack.Field{Title: "申請者", Value: strconv.Itoa(int(result.User.ID)) + ":" + result.User.Name}).
		AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(result.Group.ID)) + ":" + result.Group.Org}).
		AddField(slack.Field{Title: "接続コード（新規発番）", Value: input.ConnectionType + fmt.Sprintf("%03d", number)}).
		AddField(slack.Field{Title: "接続コード（補足情報）", Value: input.ConnectionComment})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

	c.JSON(http.StatusOK, group.Result{})
}
