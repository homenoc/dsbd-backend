package v0

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	serviceTemplate "github.com/homenoc/dsbd-backend/pkg/api/core/template/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbServiceTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/service/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func Add(c *gin.Context) {
	var input service.Input
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// group authentication
	result := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	// check user level
	if result.User.Level > 2 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	// check json
	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// status check for group
	if !(*result.User.Group.ExpiredStatus == 0 && *result.User.Group.Pass) {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: failed group status"})
		return
	}

	// add_allow check for group
	if !(*result.User.Group.AddAllow) {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: failed group add_allow status"})
		return
	}

	var grpIP []core.IP = nil

	resultServiceTemplate := dbServiceTemplate.Get(serviceTemplate.ID, &core.ServiceTemplate{Model: gorm.Model{ID: input.ServiceTemplateID}})
	if resultServiceTemplate.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultServiceTemplate.Err.Error()})
		return
	}

	if *resultServiceTemplate.Services[0].NeedJPNIC {
		if err = checkJPNIC(input); err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}

		if err = checkJPNICAdminUser(input.JPNICAdmin); err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}

		if len(input.JPNICTech) == 0 || len(input.JPNICTech) > 2 {
			c.JSON(http.StatusBadRequest, common.Error{Error: "error: user tech count"})
			return
		}

		for _, tmp := range input.JPNICTech {
			if err = checkJPNICTechUser(tmp); err != nil {
				c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
				return
			}
		}

		grpIP, err = ipProcess(false, true, input.IP)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	if *resultServiceTemplate.Services[0].NeedComment && input.ServiceComment == "" {
		c.JSON(http.StatusBadRequest, common.Error{Error: "no data: comment"})
		return
	}

	if *resultServiceTemplate.Services[0].NeedGlobalAS {
		if input.ASN == 0 {
			c.JSON(http.StatusBadRequest, common.Error{Error: "no data: ASN"})
			return
		}

		grpIP, err = ipProcess(false, false, input.IP)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	resultNetwork := dbService.Get(service.SearchNewNumber, &core.Service{GroupID: result.User.Group.ID})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultNetwork.Err.Error()})
		return
	}
	var number uint = 1
	for _, tmp := range resultNetwork.Service {
		if tmp.ServiceNumber >= 1 {
			number = tmp.ServiceNumber + 1
		}
	}

	if number >= 999 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "error: over number"})
		return
	}

	// db create for network
	net, err := dbService.Create(&core.Service{
		GroupID:           result.User.Group.ID,
		ServiceTemplateID: &input.ServiceTemplateID,
		ServiceComment:    input.ServiceComment,
		ServiceNumber:     number,
		Org:               input.Org,
		OrgEn:             input.OrgEn,
		PostCode:          input.Postcode,
		Address:           input.Address,
		AddressEn:         input.AddressEn,
		AveUpstream:       input.AveUpstream,
		MaxUpstream:       input.MaxUpstream,
		AveDownstream:     input.AveDownstream,
		MaxDownstream:     input.MaxDownstream,
		ASN:               &[]uint{input.ASN}[0],
		IP:                grpIP,
		JPNICAdmin:        input.JPNICAdmin,
		JPNICTech:         input.JPNICTech,
		Enable:            &[]bool{true}[0],
		Pass:              &[]bool{false}[0],
		AddAllow:          &[]bool{true}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "ネットワーク情報登録"}).
		AddField(slack.Field{Title: "申請者", Value: strconv.Itoa(int(result.User.ID)) + ":" + result.User.Name}).
		AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(result.User.Group.ID)) + ":" + result.User.Group.Org}).
		AddField(slack.Field{Title: "サービスコード（新規発番）", Value: resultServiceTemplate.Services[0].Type + fmt.Sprintf("%03d", number)}).
		AddField(slack.Field{Title: "サービスコード（補足情報）", Value: input.ServiceComment})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

	// ---------ここまで処理が通っている場合、DBへの書き込みにすべて成功している
	// GroupのStatusをAfterStatusにする
	if err = dbGroup.Update(group.UpdateAll, core.Group{
		Model:    gorm.Model{ID: result.User.Group.ID},
		AddAllow: &[]bool{false}[0],
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	attachment = slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "ステータス変更"}).
		AddField(slack.Field{Title: "申請者", Value: "System"}).
		AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(result.User.Group.ID)) + ":" + result.User.Group.Org}).
		AddField(slack.Field{Title: "現在ステータス情報", Value: "審査中"}).
		AddField(slack.Field{Title: "ステータス履歴", Value: "1[ネットワーク情報記入段階(User)] =>2[審査中] "})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

	c.JSON(http.StatusOK, service.ResultOne{Service: *net})
}

// Todo: 以下の処理は実装中
func Update(c *gin.Context) {
	var input core.Service
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
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

	resultNetwork := dbService.Get(service.ID, &core.Service{Model: gorm.Model{ID: input.ID}})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNetwork.Err.Error()})
		return
	}
	if len(resultNetwork.Service) == 0 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "failed Service ID"})
		return
	}
	if resultNetwork.Service[0].GroupID != result.User.Group.ID {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "Authentication failure"})
		return
	}

	replace := replaceService(resultNetwork.Service[0], input)

	if err = dbService.Update(service.UpdateData, replace); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.Result{})
}

func GetAddAllow(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	if resultService := dbService.Get(service.GIDAndAddAllow, &core.Service{GroupID: result.User.Group.ID}); resultService.Err != nil {
		log.Println(resultService.Err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultService.Err.Error()})
	} else {
		c.JSON(http.StatusOK, service.Result{Service: resultService.Service})
	}
}
