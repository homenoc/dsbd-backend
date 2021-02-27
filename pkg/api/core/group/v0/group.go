package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicTech"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbJpnicAdmin "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/jpnicAdmin/v0"
	dbJpnicTech "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/jpnicTech/v0"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

//参照関連のエラーが出る可能性あるかもしれない
func Add(c *gin.Context) {
	var input group.Group
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	userResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
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

	result, err := dbGroup.Create(&group.Group{
		Agree:     &[]bool{true}[0],
		Question:  input.Question,
		Org:       input.Org,
		Status:    &[]uint{1}[0],
		Bandwidth: input.Bandwidth,
		Comment:   input.Comment,
		Contract:  input.Contract,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "グループ登録"}).
		AddField(slack.Field{Title: "Question", Value: input.Question}).
		AddField(slack.Field{Title: "Bandwidth", Value: input.Bandwidth}).
		AddField(slack.Field{Title: "Contract", Value: input.Contract})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

	if err = dbUser.Update(user.UpdateGID, &user.User{Model: gorm.Model{ID: userResult.User.ID}, GroupID: result.Model.ID}); err != nil {
		log.Println(dbGroup.Delete(&group.Group{Model: gorm.Model{ID: result.ID}}))
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, common.Result{})
	}
}

func Update(c *gin.Context) {
	var input group.Group

	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	log.Println(c.BindJSON(&input))

	authResult := auth.GroupAuthentication(0, token.Token{UserToken: userToken, AccessToken: accessToken})
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

	if data.Bandwidth != input.Bandwidth {
		data.Bandwidth = input.Bandwidth
	}

	if err := dbGroup.Update(group.UpdateInfo, data); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: authResult.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{})

}

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(1, token.Token{UserToken: userToken, AccessToken: accessToken})
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

	resultGroup := dbGroup.Get(group.ID, &group.Group{Model: gorm.Model{ID: result.Group.ID}})
	if resultGroup.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	// Network情報にて開通しているものを抜き出す
	resultNetwork := dbNetwork.Get(network.Open, &network.Network{GroupID: result.Group.ID})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	open := false
	if len(resultNetwork.Network) > 0 {
		open = true
	}

	c.JSON(http.StatusOK, group.ResultOne{
		ID:            resultGroup.Group[0].ID,
		Agree:         resultGroup.Group[0].Agree,
		Question:      resultGroup.Group[0].Question,
		Org:           resultGroup.Group[0].Org,
		Status:        *resultGroup.Group[0].Status,
		Bandwidth:     resultGroup.Group[0].Bandwidth,
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

	result := auth.GroupAuthentication(0, token.Token{UserToken: userToken, AccessToken: accessToken})
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

	resultNetwork := dbNetwork.Get(network.GID, &network.Network{GroupID: result.Group.ID})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}
	if len(resultNetwork.Network) == 0 {
		resultNetwork.Network = nil
	}

	resultConnection := dbConnection.Get(connection.GID, &connection.Connection{GroupID: result.Group.ID})
	if resultConnection.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}
	if len(resultConnection.Connection) == 0 {
		resultConnection.Connection = nil
	}

	var resultJpnicTech []jpnicTech.JpnicTech = nil
	var resultJpnicAdmin []jpnicAdmin.JpnicAdmin = nil

	for _, data := range resultNetwork.Network {
		tmpAdmin := dbJpnicAdmin.Get(jpnicAdmin.NetworkId, &jpnicAdmin.JpnicAdmin{NetworkID: data.ID})
		if tmpAdmin.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: tmpAdmin.Err.Error()})
			return
		}
		if len(tmpAdmin.Jpnic) == 0 {
			break
		}
		resultJpnicAdmin = append(resultJpnicAdmin, tmpAdmin.Jpnic[0])

		tmpTech := dbJpnicTech.Get(jpnicAdmin.NetworkId, &jpnicTech.JpnicTech{NetworkID: data.ID})
		if tmpAdmin.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: tmpAdmin.Err.Error()})
			return
		}
		if len(tmpTech.Jpnic) == 0 {
			break
		}
		for _, tmpTechDetail := range tmpTech.Jpnic {
			resultJpnicTech = append(resultJpnicTech, tmpTechDetail)
		}
	}

	c.JSON(http.StatusOK, group.ResultAll{
		Group: group.ResultOne{
			ID:        result.Group.ID,
			Agree:     result.Group.Agree,
			Question:  result.Group.Question,
			Org:       result.Group.Org,
			Status:    *result.Group.Status,
			Bandwidth: result.Group.Bandwidth,
			Contract:  result.Group.Contract,
			Student:   result.Group.Student,
		},
		Network:    resultNetwork.Network,
		JpnicAdmin: resultJpnicAdmin,
		JpnicTech:  resultJpnicTech,
		Connection: resultConnection.Connection,
	})
}
