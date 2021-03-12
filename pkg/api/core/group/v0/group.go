package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
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
	if userResult.User.Level > 1 {
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
		Status:         &[]uint{1}[0],
		ExpiredStatus:  &[]uint{0}[0],
		Contract:       input.Contract,
		Student:        input.Student,
		StudentExpired: studentExpired,
		Lock:           &[]bool{true}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "グループ登録"}).
		AddField(slack.Field{Title: "Question", Value: input.Question}).
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

	if authResult.User.Level > 1 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: failed user level"})
		return
	}
	if *authResult.Group.Lock {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: This group is locked"})
		return
	}

	data := authResult.Group

	if data.Org != input.Org {
		data.Org = input.Org
	}

	if err = dbGroup.Update(group.UpdateInfo, data); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: authResult.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{})

}

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(1, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	if result.User.Level >= 10 {
		if result.User.Level > 1 {
			c.JSON(http.StatusInternalServerError, common.Error{Error: "You don't have authority this operation"})
			return
		}
	}

	resultGroup := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: result.Group.ID}})
	if resultGroup.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	// Network情報にて開通しているものを抜き出す
	resultService := dbService.Get(service.Open, &core.Service{GroupID: result.Group.ID})
	if resultService.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	open := false
	if len(resultService.Service) > 0 {
		open = true
	}

	c.JSON(http.StatusOK, group.ResultOne{
		ID:            resultGroup.Group[0].ID,
		Agree:         resultGroup.Group[0].Agree,
		Question:      resultGroup.Group[0].Question,
		Org:           resultGroup.Group[0].Org,
		Status:        *resultGroup.Group[0].Status,
		Contract:      resultGroup.Group[0].Contract,
		Student:       resultGroup.Group[0].Student,
		Pass:          resultGroup.Group[0].Pass,
		Lock:          resultGroup.Group[0].Lock,
		ExpiredStatus: *resultGroup.Group[0].ExpiredStatus,
		Open:          &open,
	})
}

func GetAll(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	if result.User.Level >= 10 {
		if result.User.Level > 1 {
			c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
			return
		}
	}

	//resultService := dbService.Get(service.GID, &core.Service{GroupID: result.Group.ID})
	//if resultService.Err != nil {
	//	c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	//	return
	//}
	//if len(resultService.Service) == 0 {
	//	resultService.Service = nil
	//}
	//
	//resultConnection := dbConnection.Get(connection.ServiceID, &core.Connection{GroupID: result.Group.ID})
	//if resultConnection.Err != nil {
	//	c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	//	return
	//}
	//if len(resultConnection.Connection) == 0 {
	//	resultConnection.Connection = nil
	//}
	//
	//var resultTech []core.JPNICTech = nil
	//var resultAdmin []core.JPNICAdmin = nil
	//
	//for _, data := range resultService.Service {
	//	tmpAdmin := dbAdmin.Get(admin.NetworkId, &admin.Admin{NetworkID: data.ID})
	//	if tmpAdmin.Err != nil {
	//		c.JSON(http.StatusInternalServerError, common.Error{Error: tmpAdmin.Err.Error()})
	//		return
	//	}
	//	if len(tmpAdmin.Admins) == 0 {
	//		break
	//	}
	//	resultAdmin = append(resultAdmin, tmpAdmin.Admins[0])
	//
	//	tmpTech := dbTech.Get(admin.NetworkId, &tech.Tech{NetworkID: data.ID})
	//	if tmpAdmin.Err != nil {
	//		c.JSON(http.StatusInternalServerError, common.Error{Error: tmpAdmin.Err.Error()})
	//		return
	//	}
	//	if len(tmpTech.Tech) == 0 {
	//		break
	//	}
	//	for _, tmpTechDetail := range tmpTech.Tech {
	//		resultTech = append(resultTech, tmpTechDetail)
	//	}
	//}

	c.JSON(http.StatusOK, group.Result{Group: result.Group})

	//c.JSON(http.StatusOK, group.ResultAll{
	//	Group: group.ResultOne{
	//		ID:       result.Group.ID,
	//		Agree:    result.Group.Agree,
	//		Question: result.Group.Question,
	//		Org:      result.Group.Org,
	//		Status:   *result.Group.Status,
	//		Contract: result.Group.Contract,
	//		Student:  result.Group.Student,
	//	},
	//	Network:    resultService.Service,
	//	Admin:      resultAdmin,
	//	Tech:       resultTech,
	//	Connection: resultcore.Connection,
	//})
}
