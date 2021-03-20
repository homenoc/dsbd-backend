package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

//参照関連のエラーが出る可能性あるかもしれない
func Add(c *gin.Context) {
	var input group.Input
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	userResult := auth.UserAuthentication(core.Token{UserToken: userToken, AccessToken: accessToken})
	if userResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: userResult.Err.Error()})
		return
	}

	// check authority
	if userResult.User.Level > 2 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	if userResult.User.GroupID != 0 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: You can't create new group"})
		return
	}

	if err = check(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	var studentExpired *time.Time = nil
	if *input.Student {
		tmpStudentExpired, _ := time.Parse("2006-01-02", *input.StudentExpired)
		studentExpired = &tmpStudentExpired
	}

	result, err := dbGroup.Create(&core.Group{
		Agree:          &[]bool{*input.Agree}[0],
		Question:       input.Question,
		Org:            input.Org,
		OrgEn:          input.OrgEn,
		PostCode:       input.PostCode,
		Address:        input.Address,
		AddressEn:      input.AddressEn,
		Tel:            input.Tel,
		Country:        input.Country,
		Status:         &[]uint{1}[0],
		ExpiredStatus:  &[]uint{0}[0],
		Contract:       input.Contract,
		Student:        input.Student,
		StudentExpired: studentExpired,
		Open:           &[]bool{false}[0],
		Pass:           &[]bool{false}[0],
		Lock:           &[]bool{true}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "グループ登録"}).
		AddField(slack.Field{Title: "Question", Value: input.Question}).
		AddField(slack.Field{Title: "Org", Value: input.Org + "(" + input.OrgEn + ")"}).
		AddField(slack.Field{Title: "Country", Value: input.Country}).
		AddField(slack.Field{Title: "Student", Value: strconv.FormatBool(*input.Student)}).
		AddField(slack.Field{Title: "Contract", Value: input.Contract})

	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

	if err = dbUser.Update(user.UpdateGID, &core.User{Model: gorm.Model{ID: userResult.User.ID}, GroupID: result.Model.ID}); err != nil {
		log.Println(dbGroup.Delete(&core.Group{Model: gorm.Model{ID: result.ID}}))
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, common.Result{})
	}
}

func Update(c *gin.Context) {
	var input core.Group

	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	authResult := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if authResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: authResult.Err.Error()})
		return
	}

	if authResult.User.Level > 2 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: failed user level"})
		return
	}
	if *authResult.User.Group.Lock {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: This group is locked"})
		return
	}

	data := authResult.User.Group

	if data.Org != input.Org {
		data.Org = input.Org
	}

	//if err = dbGroup.Update(group.UpdateInfo, data); err != nil {
	//	c.JSON(http.StatusInternalServerError, common.Error{Error: authResult.Err.Error()})
	//	return
	//}
	c.JSON(http.StatusOK, common.Result{})

}

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(1, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	resultUser := dbUser.Get(user.ID, &core.User{Model: gorm.Model{ID: result.User.ID}})
	if resultUser.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultUser.Err.Error()})
		return
	}

	var responseUser []user.User
	if result.User.Level == 1 || result.User.Level == 2 {
		for _, tmpUser := range resultUser.User {
			responseUser = append(responseUser, user.User{
				ID:            tmpUser.ID,
				Name:          tmpUser.Name,
				NameEn:        tmpUser.NameEn,
				Email:         tmpUser.Email,
				Level:         tmpUser.Level,
				ExpiredStatus: *tmpUser.ExpiredStatus,
				MailVerify:    tmpUser.MailVerify,
			})
		}
	} else {
		responseUser = []user.User{
			{
				ID:            result.User.ID,
				Name:          result.User.Name,
				NameEn:        result.User.NameEn,
				Email:         result.User.Email,
				Level:         result.User.Level,
				ExpiredStatus: *result.User.ExpiredStatus,
				MailVerify:    result.User.MailVerify,
			},
		}
	}

	log.Println(result.User)

	c.JSON(http.StatusOK, group.Result{Group: group.Group{
		ID:            result.User.Group.ID,
		Agree:         result.User.Group.Agree,
		Question:      result.User.Group.Question,
		Org:           result.User.Group.Org,
		Status:        *result.User.Group.Status,
		Contract:      result.User.Group.Contract,
		Student:       result.User.Group.Student,
		Open:          result.User.Group.Open,
		Pass:          result.User.Group.Pass,
		Lock:          result.User.Group.Lock,
		ExpiredStatus: *result.User.Group.ExpiredStatus,
		//User:          users,
	}})
}

func GetAll(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	if result.User.Level > 2 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	resultGroup := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: result.User.GroupID}})
	if resultGroup.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	var users []user.User
	var services []service.Service

	for _, tmpUser := range resultGroup.Group[0].Users {
		users = append(users, user.User{
			ID:            tmpUser.ID,
			Name:          tmpUser.Name,
			NameEn:        tmpUser.NameEn,
			Email:         tmpUser.Email,
			Level:         tmpUser.Level,
			ExpiredStatus: *tmpUser.ExpiredStatus,
			MailVerify:    tmpUser.MailVerify,
		})
	}

	for _, tmpService := range resultGroup.Group[0].Services {
		var tmpConnections *[]connection.Connection = nil
		var tmpJPNICAdmin *service.JPNIC = nil
		var tmpJPNICTech *[]service.JPNIC = nil

		if tmpService.JPNICAdmin.ID == 0 {
			tmpJPNICAdmin = &service.JPNIC{
				ID:          tmpService.JPNICAdmin.ID,
				JPNICHandle: tmpService.JPNICAdmin.JPNICHandle,
				Name:        tmpService.JPNICAdmin.Name,
				NameEn:      tmpService.JPNICAdmin.NameEn,
				Org:         tmpService.JPNICAdmin.Org,
				OrgEn:       tmpService.JPNICAdmin.OrgEn,
				PostCode:    tmpService.JPNICAdmin.PostCode,
				Address:     tmpService.JPNICAdmin.Address,
				AddressEn:   tmpService.JPNICAdmin.AddressEn,
				Dept:        tmpService.JPNICAdmin.Dept,
				DeptEn:      tmpService.JPNICAdmin.DeptEn,
				Tel:         tmpService.JPNICAdmin.Tel,
				Fax:         tmpService.JPNICAdmin.Fax,
				Country:     tmpService.JPNICAdmin.Country,
			}
		}

		if tmpService.JPNICTech != nil || len(tmpService.JPNICTech) != 0 {
			for _, tmpServiceJPNICTech := range tmpService.JPNICTech {
				*tmpJPNICTech = append(*tmpJPNICTech, service.JPNIC{
					ID:          tmpServiceJPNICTech.ID,
					JPNICHandle: tmpServiceJPNICTech.JPNICHandle,
					Name:        tmpServiceJPNICTech.Name,
					NameEn:      tmpServiceJPNICTech.NameEn,
					Org:         tmpServiceJPNICTech.Org,
					OrgEn:       tmpServiceJPNICTech.OrgEn,
					PostCode:    tmpServiceJPNICTech.PostCode,
					Address:     tmpServiceJPNICTech.Address,
					AddressEn:   tmpServiceJPNICTech.AddressEn,
					Dept:        tmpServiceJPNICTech.Dept,
					DeptEn:      tmpServiceJPNICTech.DeptEn,
					Tel:         tmpServiceJPNICTech.Tel,
					Fax:         tmpServiceJPNICTech.Fax,
					Country:     tmpServiceJPNICTech.Country,
				})
			}
		}

		if tmpConnections != nil {
			for _, tmpConnection := range tmpService.Connection {
				*tmpConnections = append(*tmpConnections, connection.Connection{
					ID:                           tmpConnection.ID,
					BGPRouterID:                  tmpConnection.BGPRouterID,
					BGPRouterName:                tmpConnection.BGPRouter.HostName,
					TunnelEndPointRouterIPID:     tmpConnection.TunnelEndPointRouterIPID,
					TunnelEndPointRouterIPIDName: tmpConnection.TunnelEndPointRouterIP.TunnelEndPointRouter.HostName,
					ConnectionTemplateID:         tmpConnection.ConnectionTemplateID,
					ConnectionTemplateName:       tmpConnection.ConnectionTemplate.Name,
					ConnectionComment:            tmpConnection.ConnectionComment,
					ConnectionNumber:             tmpConnection.ConnectionNumber,
					NTTTemplateID:                tmpConnection.NTTTemplateID,
					NTTTemplateName:              tmpConnection.NTTTemplate.Name,
					NOCID:                        tmpConnection.NOCID,
					NOCName:                      tmpConnection.NOC.Name,
					TermIP:                       tmpConnection.TermIP,
					Monitor:                      tmpConnection.Monitor,
					Address:                      tmpConnection.Address,
					LinkV4Our:                    tmpConnection.LinkV4Our,
					LinkV4Your:                   tmpConnection.LinkV4Your,
					LinkV6Our:                    tmpConnection.LinkV6Our,
					LinkV6Your:                   tmpConnection.LinkV6Your,
					Open:                         tmpConnection.Open,
					Lock:                         tmpConnection.Lock,
				})
			}
		}

		services = append(services, service.Service{
			ID:                  tmpService.ID,
			GroupID:             tmpService.ID,
			ServiceTemplateID:   tmpService.ServiceTemplateID,
			ServiceTemplateName: tmpService.ServiceTemplate.Name,
			ServiceComment:      tmpService.ServiceComment,
			ServiceNumber:       tmpService.ServiceNumber,
			Org:                 tmpService.Org,
			OrgEn:               tmpService.OrgEn,
			PostCode:            tmpService.PostCode,
			Address:             tmpService.Address,
			AddressEn:           tmpService.AddressEn,
			ASN:                 tmpService.ASN,
			RouteV4:             tmpService.RouteV4,
			RouteV6:             tmpService.RouteV6,
			V4Name:              tmpService.V4Name,
			V6Name:              tmpService.V6Name,
			AveUpstream:         tmpService.AveUpstream,
			MaxUpstream:         tmpService.MaxUpstream,
			AveDownstream:       tmpService.AveDownstream,
			MaxDownstream:       tmpService.MaxDownstream,
			MaxBandWidthAS:      tmpService.MaxBandWidthAS,
			Fee:                 tmpService.Fee,
			IP:                  tmpService.IP,
			Connections:         tmpConnections,
			JPNICAdminID:        tmpService.JPNICAdminID,
			JPNICAdmin:          tmpJPNICAdmin,
			JPNICTech:           tmpJPNICTech,
			Open:                tmpService.Open,
			AddAllow:            tmpService.AddAllow,
			Lock:                tmpService.Lock,
		})
	}

	c.JSON(http.StatusOK, group.Result{Group: group.Group{
		ID:            resultGroup.Group[0].ID,
		Agree:         resultGroup.Group[0].Agree,
		Question:      resultGroup.Group[0].Question,
		Org:           resultGroup.Group[0].Org,
		Status:        *resultGroup.Group[0].Status,
		Contract:      resultGroup.Group[0].Contract,
		Student:       resultGroup.Group[0].Student,
		Open:          resultGroup.Group[0].Open,
		Pass:          resultGroup.Group[0].Pass,
		Lock:          resultGroup.Group[0].Lock,
		ExpiredStatus: *resultGroup.Group[0].ExpiredStatus,
		User:          users,
		Service:       &services,
	}})
}
