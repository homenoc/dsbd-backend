package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	groupConnection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	groupNetwork "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input group.Group

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if _, err = dbGroup.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.Result{})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err := dbGroup.Delete(&group.Group{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.Result{})
}

func UpdateAdmin(c *gin.Context) {
	var input group.Group

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	tmp := dbGroup.Get(group.ID, &group.Group{Model: gorm.Model{ID: input.ID}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: tmp.Err.Error()})
		return
	}

	// 審査ステータスのSlack通知
	if *tmp.Group[0].Pass != *input.Pass {
		attachment := slack.Attachment{}
		if *input.Pass {
			attachment.AddField(slack.Field{Title: "Title", Value: "ステータス変更"}).
				AddField(slack.Field{Title: "申請者", Value: "管理者"}).
				AddField(slack.Field{Title: "Group", Value: strconv.Itoa(int(tmp.Group[0].ID)) + ":" + tmp.Group[0].Org}).
				AddField(slack.Field{Title: "現在ステータス情報", Value: "審査中"}).
				AddField(slack.Field{Title: "ステータス履歴", Value: "[審査合格] =>[審査中] "})
			notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

		} else if !(*input.Pass) {
			attachment.AddField(slack.Field{Title: "Title", Value: "ステータス変更"}).
				AddField(slack.Field{Title: "申請者", Value: "管理者"}).
				AddField(slack.Field{Title: "Group", Value: strconv.Itoa(int(tmp.Group[0].ID)) + ":" + tmp.Group[0].Org}).
				AddField(slack.Field{Title: "現在ステータス情報", Value: "審査合格処理終了"}).
				AddField(slack.Field{Title: "ステータス履歴", Value: "[審査中] =>[審査合格処理]"})
			notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

		}
	}

	replace, err := updateAdminGroup(input, tmp.Group[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: this email is already registered"})
		return
	}

	if err = dbGroup.Update(group.UpdateAll, replace); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.Result{})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := dbGroup.Get(group.ID, &group.Group{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	resultUser := dbUser.Get(user.GID, &user.User{GroupID: uint(id)})
	if resultUser.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	resultNetwork := dbNetwork.Get(groupNetwork.GID, &groupNetwork.Network{GroupID: uint(id)})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	resultConnection := dbConnection.Get(groupConnection.GID, &groupConnection.Connection{GroupID: uint(id)})
	if resultConnection.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, group.AdminResult{User: resultUser.User, Group: result.Group,
		Network: resultNetwork.Network, Connection: resultConnection.Connection})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbGroup.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, group.Result{Group: result.Group})
	}
}
