package v0

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/ip"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/jpnicTech"
	serviceTemplate "github.com/homenoc/dsbd-backend/pkg/api/core/template/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbPlan "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/ip/plan/v0"
	dbIP "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/ip/v0"
	dbJPNICTech "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/jpnicTech/v0"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	dbServiceTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/service/v0"
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

	// serviceIDが0の時エラー処理
	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("This id is wrong... ")})
		return
	}

	var input service.Input

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
	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	log.Println(input)

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

		grpIP, err = ipProcess(true, true, input.IP)
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

		grpIP, err = ipProcess(true, false, input.IP)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	if *resultServiceTemplate.Services[0].NeedRoute && input.RouteV4 == "" && input.RouteV6 == "" {
		c.JSON(http.StatusBadRequest, common.Error{Error: "no data: Route Information"})
		return
	}

	resultNetwork := dbService.Get(service.SearchNewNumber, &core.Service{GroupID: uint(id)})
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

	//db create for network
	_, err = dbService.Create(&core.Service{
		GroupID:           uint(id),
		ServiceTemplateID: &input.ServiceTemplateID,
		ServiceComment:    input.ServiceComment,
		ServiceNumber:     number,
		Org:               input.Org,
		OrgEn:             input.OrgEn,
		PostCode:          input.Postcode,
		Address:           input.Address,
		AddressEn:         input.AddressEn,
		RouteV4:           input.RouteV4,
		RouteV6:           input.RouteV6,
		AveUpstream:       input.AveUpstream,
		MaxUpstream:       input.MaxUpstream,
		AveDownstream:     input.AveDownstream,
		MaxDownstream:     input.MaxDownstream,
		MaxBandWidthAS:    input.MaxBandWidthAS,
		ASN:               &[]uint{input.ASN}[0],
		Fee:               &[]uint{0}[0],
		IP:                grpIP,
		JPNICAdmin:        input.JPNICAdmin,
		JPNICTech:         input.JPNICTech,
		Open:              &[]bool{false}[0],
		Lock:              &[]bool{true}[0],
		AddAllow:          &[]bool{false}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "ネットワーク情報登録(管理者実行)"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(id)}).
		AddField(slack.Field{Title: "サービスコード（新規発番）", Value: resultServiceTemplate.Services[0].Type + fmt.Sprintf("%03d", number)}).
		AddField(slack.Field{Title: "サービスコード（補足情報）", Value: input.ServiceComment})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

	c.JSON(http.StatusOK, service.Result{})
}

func AddIPAdmin(c *gin.Context) {
	var input service.IPInput

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

	resultIP, err := ipProcess(true, false, []service.IPInput{input})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.JoinIP(uint(id), resultIP[0]); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, service.Result{})
}

func AddPlanAdmin(c *gin.Context) {
	var input core.Plan

	err := c.BindJSON(&input)
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

	if _, err = dbPlan.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, service.Result{})
}

func AddJPNICAdminAdmin(c *gin.Context) {
	var input core.JPNICAdmin

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

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.JoinJPNICAdmin(uint(id), input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{})
}

func AddJPNICTechAdmin(c *gin.Context) {
	var input core.JPNICTech

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

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.JoinJPNICTech(uint(id), input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{})
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

	if err := dbService.Delete(&core.Service{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, service.Result{})
}

func DeleteIPAdmin(c *gin.Context) {
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

	ipID, err := strconv.Atoi(c.Param("ip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.DeleteIP(uint(id), uint(ipID)); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{})
}

func DeletePlanAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	planID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbPlan.Delete(&core.Plan{Model: gorm.Model{ID: uint(planID)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{})
}

func DeleteJPNICAdminAdmin(c *gin.Context) {
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

	jpnicID, err := strconv.Atoi(c.Param("jpnic_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.DeleteJPNICAdmin(uint(id), uint(jpnicID)); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{})
}

func DeleteJPNICTechAdmin(c *gin.Context) {
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

	jpnicID, err := strconv.Atoi(c.Param("jpnic_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.DeleteJPNICTech(uint(id), uint(jpnicID)); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{})
}

func UpdateAdmin(c *gin.Context) {
	var input core.Service

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

	resultService := dbService.Get(service.ID, &core.Service{Model: gorm.Model{ID: uint(id)}})
	if resultService.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultService.Err.Error()})
		return
	}

	noticeSlackAdmin(resultService.Service[0], input)

	input.ID = uint(id)

	if err = dbService.Update(service.UpdateAll, input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, service.Result{})
}

func UpdateIPAdmin(c *gin.Context) {
	var input service.IPInput

	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	ipID, err := strconv.Atoi(c.Param("ip_id"))
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

	resultIP, err := ipProcess(true, false, []service.IPInput{input})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}
	resultIP[0].ID = uint(ipID)
	log.Println(input)

	if err = dbIP.Update(ip.UpdateAll, resultIP[0]); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, service.Result{})
}

func UpdatePlanAdmin(c *gin.Context) {
	var input core.Plan

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

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	input.ID = uint(id)

	if err = dbPlan.Update(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{})
}

func UpdateJPNICAdminAdmin(c *gin.Context) {
	var input core.JPNICAdmin

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

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.UpdateJPNICAdmin(uint(id), input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{})
}

func UpdateJPNICTechAdmin(c *gin.Context) {
	var input core.JPNICTech

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

	jpnicID, err := strconv.Atoi(c.Param("jpnic_id"))
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

	ok := false

	resultJPNICTech := dbService.Get(service.ID, &core.Service{Model: gorm.Model{ID: uint(id)}})
	if resultJPNICTech.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultJPNICTech.Err.Error()})
		return
	}

	for _, tmp := range resultJPNICTech.Service[0].JPNICTech {
		if tmp.ID == uint(jpnicID) {
			ok = true
			input.Model = tmp.Model
			break
		}
	}

	if ok {
		if err = dbJPNICTech.Update(jpnicTech.UpdateAll, input); err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.Result{})
	} else {
		c.JSON(http.StatusBadRequest, common.Error{Error: "Not Found GroupID"})
	}
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

	result := dbService.Get(service.ID, &core.Service{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, service.Result{Service: result.Service})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbService.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, service.Result{Service: result.Service})
	}
}
