package v0

import (
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

	// check lock
	if *result.Group.Lock {
		c.JSON(http.StatusForbidden, common.Error{Error: "Lock status"})
		return
	}

	// check json
	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, group.Result{Error: err.Error()})
		return
	}

	// status check for group
	if !(*result.Group.Status == 1 && *result.Group.ExpiredStatus == 0 && *result.Group.Pass) {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: failed group status"})
		return
	}

	grpIP, err := ipProcess(input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	jh := jpnicHandler{
		admin: input.AdminID, tech: input.TechID, groupID: result.Group.ID, jpnicAdmin: nil, jpnicTech: nil,
	}

	// PIアドレスではない場合、jpnic Processを実行
	if !input.PI {
		if err = jh.jpnicProcess(); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	// db create for network
	net, err := dbNetwork.Create(&network.Network{
		GroupID:    result.Group.ID,
		Org:        input.Org,
		OrgEn:      input.OrgEn,
		Postcode:   input.Postcode,
		Address:    input.Address,
		AddressEn:  input.AddressEn,
		RouteV4:    input.RouteV4,
		RouteV6:    input.RouteV6,
		PI:         &[]bool{input.PI}[0],
		ASN:        input.ASN,
		Open:       &[]bool{false}[0],
		IP:         *grpIP,
		JPNICAdmin: *jh.jpnicAdmin,
		JPNICTech:  *jh.jpnicTech,
		Lock:       &[]bool{input.Lock}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "ネットワーク登録"}).
		AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(input.GroupID))})
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

	log.Println(c.BindJSON(&input))

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

	//if !(result.Group.Status == 211 || result.Group.Status == 221) {
	//	c.JSON(http.StatusUnauthorized, common.Error{Error: "error: group status"})
	//	return
	//}

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
