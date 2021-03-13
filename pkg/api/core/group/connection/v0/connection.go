package v0

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	connectionTemplate "github.com/homenoc/dsbd-backend/pkg/api/core/template/connection"
	ntt "github.com/homenoc/dsbd-backend/pkg/api/core/template/ntt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	dbConnectionTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/connection/v0"
	dbNTTTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/ntt/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func Add(c *gin.Context) {
	var input connection.Input
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// ID取得
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: ID is 0"})
		return
	}

	result := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	// check authority
	if result.User.Level > 2 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	// status check for group
	if !*result.Group.Pass {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: Your group has not yet been reviewed."})
		return
	}

	if !(*result.Group.Status == 3 && *result.Group.ExpiredStatus == 0) {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: failed group status"})
		return
	}

	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultConnectionTemplate := dbConnectionTemplate.Get(connectionTemplate.ID,
		&core.ConnectionTemplate{Model: gorm.Model{ID: *input.ConnectionTemplateID}})
	if resultConnectionTemplate.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultConnectionTemplate.Err.Error()})
		return
	}

	if *resultConnectionTemplate.Connections[0].NeedInternet {
		resultNTT := dbNTTTemplate.Get(ntt.ID, &core.NTTTemplate{Model: gorm.Model{ID: *input.NTTTemplateID}})
		if resultNTT.Err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: resultNTT.Err.Error()})
			return
		}
	}

	// NOCIDが0の時、「どこでも収容」という意味
	if *input.NOCID != 0 {
		resultNOC := dbNOC.Get(noc.ID, &core.NOC{Model: gorm.Model{ID: *input.NOCID}})
		if resultNOC.Err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: resultNOC.Err.Error()})
			return
		}
	}

	resultService := dbService.Get(service.IDOnlySingle, &core.Service{Model: gorm.Model{ID: uint(id)}})
	if resultService.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultService.Err.Error()})
		return
	}

	if !(*resultService.Service[0].AddAllow) {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: You are not allowed to add any connection information."})
		return
	}

	// GroupIDが一致しない場合はエラーを返す
	if resultService.Service[0].GroupID != result.Group.ID {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: GroupID does not match."})
		return
	}

	resultConnection := dbConnection.Get(connection.ServiceID, &core.Connection{ServiceID: uint(id)})
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

	_, err = dbConnection.Create(&core.Connection{
		ServiceID:            resultService.Service[0].ID,
		ConnectionTemplateID: input.ConnectionTemplateID,
		ConnectionComment:    input.ConnectionComment,
		ConnectionNumber:     number,
		NTTTemplateID:        input.NTTTemplateID,
		NOCID:                input.NOCID,
		TermIP:               input.TermIP,
		Address:              input.Address,
		Monitor:              input.Monitor,
		Open:                 &[]bool{false}[0],
		Lock:                 &[]bool{true}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "接続情報登録"}).
		AddField(slack.Field{Title: "申請者", Value: strconv.Itoa(int(result.User.ID)) + ":" + result.User.Name}).
		AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(result.Group.ID)) + ":" + result.Group.Org}).
		AddField(slack.Field{Title: "サービスコード", Value: resultService.Service[0].ServiceTemplate.Type +
			fmt.Sprintf("%03d", resultService.Service[0].ServiceNumber)}).
		AddField(slack.Field{Title: "接続コード（新規発番）", Value: resultConnectionTemplate.Connections[0].Type +
			fmt.Sprintf("%03d", number)}).
		AddField(slack.Field{Title: "接続コード（補足情報）", Value: input.ConnectionComment})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

	if err = dbGroup.Update(group.UpdateStatus, core.Group{
		Model:  gorm.Model{ID: result.Group.ID},
		Status: &[]uint{4}[0],
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.Update(service.UpdateStatus, core.Service{
		Model:    gorm.Model{ID: resultService.Service[0].ID},
		AddAllow: &[]bool{false}[0],
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	attachment = slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "ステータス変更"}).
		AddField(slack.Field{Title: "申請者", Value: "System"}).
		AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(result.Group.ID)) + ":" + result.Group.Org}).
		AddField(slack.Field{Title: "現在ステータス情報", Value: "開通作業中"}).
		AddField(slack.Field{Title: "ステータス履歴", Value: "3[接続情報記入段階(User)] =>4[開通作業中] "})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

	c.JSON(http.StatusOK, group.Result{})
}
