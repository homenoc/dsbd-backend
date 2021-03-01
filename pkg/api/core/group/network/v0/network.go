package v0

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func Add(c *gin.Context) {
	var input network.Input
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// group authentication
	result := auth.GroupAuthentication(0, token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	// check user level
	if result.User.Level > 1 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	// check json
	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, group.Result{Error: err.Error()})
		return
	}

	// status check for group
	if !(*result.Group.Status == 1 && *result.Group.ExpiredStatus == 0 && *result.Group.Pass) {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: failed group status"})
		return
	}

	var grpIP *[]network.IP

	if !(input.NetworkType == "ET00") {
		grpIP, err = ipProcess(input)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	} else {
		grpIP = nil
	}

	jh := jpnicHandler{
		admin:      input.AdminID,
		tech:       input.TechID,
		groupID:    result.Group.ID,
		jpnicAdmin: nil,
		jpnicTech:  nil,
	}

	// 2000,3S00,3B00の場合
	if input.NetworkType == "2000" || input.NetworkType == "3S00" ||
		input.NetworkType == "3B00" || input.NetworkType == "IP3B" {
		if err = jh.jpnicProcess(); err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	resultNetwork := dbNetwork.Get(network.SearchNewNumber, &network.Network{GroupID: result.Group.ID})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultNetwork.Err.Error()})
		return
	}
	var number uint = 1
	for _, tmp := range resultNetwork.Network {
		if tmp.NetworkNumber >= 1 {
			number = tmp.NetworkNumber + 1
		}
	}

	if number >= 999 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "error: over number"})
		return
	}

	// db create for network
	net, err := dbNetwork.Create(&network.Network{
		GroupID:        result.Group.ID,
		NetworkType:    input.NetworkType,
		NetworkComment: input.NetworkComment,
		NetworkNumber:  number,
		Org:            input.Org,
		OrgEn:          input.OrgEn,
		Postcode:       input.Postcode,
		Address:        input.Address,
		AddressEn:      input.AddressEn,
		RouteV4:        input.RouteV4,
		RouteV6:        input.RouteV6,
		ASN:            input.ASN,
		Open:           &[]bool{false}[0],
		IP:             *grpIP,
		JPNICAdmin:     *jh.jpnicAdmin,
		JPNICTech:      *jh.jpnicTech,
		Lock:           &[]bool{true}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "ネットワーク情報登録"}).
		AddField(slack.Field{Title: "申請者", Value: strconv.Itoa(int(result.User.ID)) + ":" + result.User.Name}).
		AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(result.Group.ID)) + ":" + result.Group.Org}).
		AddField(slack.Field{Title: "サービスコード（新規発番）", Value: input.NetworkType + fmt.Sprintf("%03d", number)}).
		AddField(slack.Field{Title: "サービスコード（補足情報）", Value: input.NetworkComment})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

	// ---------ここまで処理が通っている場合、DBへの書き込みにすべて成功している
	// GroupのStatusをAfterStatusにする
	if err = dbGroup.Update(group.UpdateStatus, group.Group{Model: gorm.Model{ID: result.Group.ID},
		Status: &[]uint{2}[0]}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, network.ResultOne{Network: *net})
}

// Todo: 以下の処理は実装中
func Update(c *gin.Context) {
	var input network.Network
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

	resultNetwork := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: input.ID}})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNetwork.Err.Error()})
		return
	}
	if len(resultNetwork.Network) == 0 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "failed Network ID"})
		return
	}
	if resultNetwork.Network[0].GroupID != result.Group.ID {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "Authentication failure"})
		return
	}
	if *resultNetwork.Network[0].Lock {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "this network is locked..."})
		return
	}

	replace := replaceNetwork(resultNetwork.Network[0], input)

	if err := dbNetwork.Update(network.UpdateData, replace); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Error: err.Error()})
	}

	c.JSON(http.StatusOK, group.Result{})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbNetwork.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, network.Result{Network: result.Network})
	}
}
