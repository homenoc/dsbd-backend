package v0

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/gen"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/hash"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/mail"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	toolToken "github.com/homenoc/dsbd-backend/pkg/api/core/tool/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Add(c *gin.Context) {
	var input, data user.User
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("wrong email address")})
		return
	}

	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	mailToken, _ := toolToken.Generate(4)

	pass := ""

	// 新規ユーザ
	if input.GroupID == 0 { //new user
		if input.Pass == "" {
			c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("wrong pass")})
			return
		}
		data = user.User{GroupID: 0, Name: input.Name, NameEn: input.NameEn, Email: input.Email, Pass: input.Pass,
			Status: 0, Level: 1, MailVerify: &[]bool{false}[0], MailToken: mailToken, GroupHandle: &[]bool{false}[0]}
		// グループ所属ユーザの登録
	} else {
		if input.Level == 0 || input.Level > 5 {
			c.JSON(http.StatusInternalServerError, common.Error{Error: fmt.Sprintf("wrong user level")})
			return
		}
		authResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
		if authResult.Err != nil {
			c.JSON(http.StatusForbidden, common.Error{Error: authResult.Err.Error()})
			return
		}
		if authResult.User.GroupID != input.GroupID && authResult.User.GroupID > 0 {
			c.JSON(http.StatusInternalServerError, common.Error{Error: "gid mismatch"})
			return
		}

		pass = gen.GenerateUUID()
		log.Println("Email: " + input.Email)
		log.Println("tmp_Pass: " + pass)

		data = user.User{GroupID: input.GroupID, Name: input.Name, NameEn: input.NameEn,
			Email: input.Email, Pass: strings.ToLower(hash.Generate(pass)), GroupHandle: input.GroupHandle,
			Status: 0, Tech: input.Tech, Level: input.Level, MailVerify: &[]bool{false}[0], MailToken: mailToken}
	}

	//check exist for database
	err = dbUser.Create(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	}
	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "ユーザ登録"}).
		AddField(slack.Field{Title: "E-Mail", Value: input.Email}).
		AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(input.GroupID))}).
		AddField(slack.Field{Title: "Name", Value: input.Name}).
		AddField(slack.Field{Title: "Name(English)", Value: input.NameEn})

	notification.SendSlack(notification.Slack{Attachment: attachment, Channel: "user", Status: true})

	if pass == "" {
		mail.SendMail(mail.Mail{
			ToMail:  data.Email,
			Subject: "本人確認のメールにつきまして",
			Content: " " + input.Name + "様\n\n" + "以下のリンクから本人確認を完了してください。\n" +
				config.Conf.Controller.User.Url + "/api/v1/verify/" + mailToken + "\n" +
				"本人確認が完了次第、ログイン可能になります。\n",
		})
	} else {
		mail.SendMail(mail.Mail{
			ToMail:  data.Email,
			Subject: "本人確認メールにつきまして",
			Content: " " + input.Name + "様\n\n" + "以下のリンクから本人確認を完了してください。\n" +
				config.Conf.Controller.User.Url + "/api/v1/verify/" + mailToken + "\n" +
				"本人確認が完了次第、ログイン可能になります。\n" + "仮パスワード: " + pass,
		})
	}

	c.JSON(http.StatusOK, user.Result{})
}

func MailVerify(c *gin.Context) {
	mailToken := c.Param("token")

	result := dbUser.Get(user.MailToken, &user.User{MailToken: mailToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error() + "| we can't find token data"})
		return
	}

	if *result.User[0].MailVerify {
		c.JSON(http.StatusInternalServerError, common.Error{Error: fmt.Sprintf("This email has already been checked")})
		return
	}
	if result.User[0].Status >= 100 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: fmt.Sprintf("error: user status")})
		return
	}

	if err := dbUser.Update(user.UpdateVerifyMail, &user.User{Model: gorm.Model{ID: result.User[0].ID},
		MailVerify: &[]bool{true}[0]}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, &user.Result{})
	}
}

func Update(c *gin.Context) {
	var input user.User

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("id error")})
		return
	}
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	log.Println(c.BindJSON(&input))

	authResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if authResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: authResult.Err.Error()})
		return
	}

	if !*authResult.User.MailVerify {
		c.JSON(http.StatusBadRequest, common.Error{Error: "not verify for user mail"})
		return
	}

	var u, serverData user.User

	if authResult.User.ID == uint(id) || id == 0 {
		serverData = authResult.User
		u.Model.ID = authResult.User.ID
		u.Status = authResult.User.Status
	} else {
		if authResult.User.GroupID == 0 {
			c.JSON(http.StatusInternalServerError, common.Error{Error: "error: Group ID = 0"})
			return
		}
		if authResult.User.Level > 1 {
			c.JSON(http.StatusInternalServerError, common.Error{Error: "error: failed user level"})
			return
		}
		userResult := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: uint(id)}})
		if userResult.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: userResult.Err.Error()})
			return
		}
		if userResult.User[0].GroupID != authResult.User.GroupID {
			c.JSON(http.StatusInternalServerError, common.Error{Error: fmt.Sprintf("failed group authentication")})
			return
		}
		serverData = userResult.User[0]
		u.Model.ID = uint(id)
		u.Status = userResult.User[0].Status
	}

	u, err = replaceUser(serverData, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if err := dbUser.Update(user.UpdateAll, &u); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, user.Result{})
	}
}

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("id error")})
		return
	}

	authResult := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	authResult.User.Pass = ""
	authResult.User.MailToken = ""
	if authResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: authResult.Err.Error()})
		return
	}

	if authResult.User.Level >= 2 {
		c.JSON(http.StatusForbidden, common.Error{Error: "You don't have the authority."})
		return
	}

	resultUser := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: uint(id)}})
	if resultUser.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultUser.Err.Error()})
		return
	}

	if resultUser.User[0].GroupID != authResult.Group.ID {
		c.JSON(http.StatusBadRequest, common.Error{Error: "GroupID is not match."})
		return
	}

	c.JSON(http.StatusOK, user.ResultOne{
		ID:          resultUser.User[0].ID,
		GroupID:     resultUser.User[0].GroupID,
		Tech:        resultUser.User[0].Tech,
		GroupHandle: resultUser.User[0].GroupHandle,
		Name:        resultUser.User[0].Name,
		NameEn:      resultUser.User[0].NameEn,
		Email:       resultUser.User[0].Email,
		Status:      resultUser.User[0].Status,
		Level:       resultUser.User[0].Level,
		MailVerify:  resultUser.User[0].MailVerify,
		Org:         resultUser.User[0].Org,
		OrgEn:       resultUser.User[0].OrgEn,
		PostCode:    resultUser.User[0].PostCode,
		Address:     resultUser.User[0].Address,
		AddressEn:   resultUser.User[0].AddressEn,
		Dept:        resultUser.User[0].Dept,
		DeptEn:      resultUser.User[0].DeptEn,
		Pos:         resultUser.User[0].Pos,
		PosEn:       resultUser.User[0].PosEn,
		Tel:         resultUser.User[0].Tel,
		Fax:         resultUser.User[0].Fax,
		Country:     resultUser.User[0].Country,
	})
}

func GetOwn(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	authResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	authResult.User.Pass = ""
	authResult.User.MailToken = ""
	if authResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: authResult.Err.Error()})
	} else {
		c.JSON(http.StatusOK, user.ResultOne{
			ID:          authResult.User.ID,
			GroupID:     authResult.User.GroupID,
			Tech:        authResult.User.Tech,
			GroupHandle: authResult.User.GroupHandle,
			Name:        authResult.User.Name,
			NameEn:      authResult.User.NameEn,
			Email:       authResult.User.Email,
			Status:      authResult.User.Status,
			Level:       authResult.User.Level,
			MailVerify:  authResult.User.MailVerify,
			Org:         authResult.User.Org,
			OrgEn:       authResult.User.OrgEn,
			PostCode:    authResult.User.PostCode,
			Address:     authResult.User.Address,
			AddressEn:   authResult.User.AddressEn,
			Dept:        authResult.User.Dept,
			DeptEn:      authResult.User.DeptEn,
			Pos:         authResult.User.Pos,
			PosEn:       authResult.User.PosEn,
			Tel:         authResult.User.Tel,
			Fax:         authResult.User.Fax,
			Country:     authResult.User.Country,
		})
	}
}

func GetGroup(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	authResult := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	result := dbUser.Get(user.GID, &user.User{GroupID: authResult.Group.ID})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	var data user.Result

	if authResult.User.Level > 1 {
		data.User = append(data.User, user.ResultOne{
			ID:          authResult.User.ID,
			GroupID:     authResult.User.GroupID,
			Tech:        authResult.User.Tech,
			GroupHandle: authResult.User.GroupHandle,
			Name:        authResult.User.Name,
			NameEn:      authResult.User.NameEn,
			Email:       authResult.User.Email,
			Status:      authResult.User.Status,
			Level:       authResult.User.Level,
			MailVerify:  authResult.User.MailVerify,
			Org:         authResult.User.Org,
			OrgEn:       authResult.User.OrgEn,
			PostCode:    authResult.User.PostCode,
			Address:     authResult.User.Address,
			AddressEn:   authResult.User.AddressEn,
			Dept:        authResult.User.Dept,
			DeptEn:      authResult.User.DeptEn,
			Pos:         authResult.User.Pos,
			PosEn:       authResult.User.PosEn,
			Tel:         authResult.User.Tel,
			Fax:         authResult.User.Fax,
			Country:     authResult.User.Country,
		})
	} else {
		resultUser := dbUser.Get(user.GID, &user.User{GroupID: authResult.Group.ID})
		if result.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
			return
		}

		for _, grp := range resultUser.User {
			data.User = append(data.User, user.ResultOne{
				ID:          grp.ID,
				GroupID:     grp.GroupID,
				Tech:        grp.Tech,
				GroupHandle: grp.GroupHandle,
				Name:        grp.Name,
				NameEn:      grp.NameEn,
				Email:       grp.Email,
				Status:      grp.Status,
				Level:       grp.Level,
				MailVerify:  grp.MailVerify,
				Org:         grp.Org,
				OrgEn:       grp.OrgEn,
				PostCode:    grp.PostCode,
				Address:     grp.Address,
				AddressEn:   grp.AddressEn,
				Dept:        grp.Dept,
				DeptEn:      grp.DeptEn,
				Pos:         grp.Pos,
				PosEn:       grp.PosEn,
				Tel:         grp.Tel,
				Fax:         grp.Fax,
				Country:     grp.Country,
			})
		}
	}

	log.Println(result)

	//for _, tmp := range result.User{
	//	tmp.Pass = ""
	//	tmp.MailToken = ""
	//	if 0 < tmp.Status && tmp.Status < 100 {
	//		data = append(data, tmp)
	//	}
	//}

	c.JSON(http.StatusOK, data)
}
