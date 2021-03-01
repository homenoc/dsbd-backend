package v0

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	// ID取得
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// networkIDが0の時エラー処理
	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("This id is wrong... ")})
		return
	}

	var input network.Input

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// check json
	if err = checkAdmin(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	log.Println(input)

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
		groupID:    uint(id),
		jpnicAdmin: nil,
		jpnicTech:  nil,
	}

	// 2000,3S00,3B00の場合
	if input.NetworkType == "2000" || input.NetworkType == "3S00" ||
		input.NetworkType == "3B00" || input.NetworkType == "CL20" || input.NetworkType == "CL3S" ||
		input.NetworkType == "CL3B" || input.NetworkType == "IP3B" {
		if err = jh.jpnicProcess(); err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	resultNetwork := dbNetwork.Get(network.SearchNewNumber, &network.Network{GroupID: uint(id)})
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

	//db create for network
	_, err = dbNetwork.Create(&network.Network{
		GroupID:        uint(id),
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
	attachment.AddField(slack.Field{Title: "Title", Value: "ネットワーク情報登録(管理者実行)"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(id)}).
		AddField(slack.Field{Title: "サービスコード（新規発番）", Value: input.NetworkType + fmt.Sprintf("%03d", number)}).
		AddField(slack.Field{Title: "サービスコード（補足情報）", Value: input.NetworkComment})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

	c.JSON(http.StatusOK, network.Result{})
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

	if err := dbNetwork.Delete(&network.Network{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, network.Result{})
}

func UpdateAdmin(c *gin.Context) {
	var input network.Network

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	result := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if err := dbNetwork.Update(network.UpdateAll, replaceAdminNetwork(result.Network[0], input)); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, network.Result{})
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

	result := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	resultUser := dbUser.Get(user.GID, &user.User{GroupID: result.Network[0].GroupID})
	if resultUser.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultUser.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, network.Result{User: resultUser.User, Network: result.Network})
}

func Get(c *gin.Context) {
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
