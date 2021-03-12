package v0

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	authInterface "github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
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
	var input user.Input
	var data core.User
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
	var authResult authInterface.GroupResult

	// 新規ユーザ
	if userToken == "" && accessToken == "" {
		if input.Pass == "" {
			c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("wrong pass")})
			return
		}

		data = core.User{
			GroupID:       0,
			Name:          input.Name,
			NameEn:        input.NameEn,
			Email:         input.Email,
			Pass:          input.Pass,
			ExpiredStatus: &[]uint{0}[0],
			Level:         1,
			MailVerify:    &[]bool{false}[0],
			MailToken:     mailToken,
		}
	} else {
		// グループ所属ユーザの登録
		resultAuth := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
		if resultAuth.Err != nil {
			c.JSON(http.StatusUnauthorized, common.Error{Error: resultAuth.Err.Error()})
			return
		}

		if input.Level > 2 {
			c.JSON(http.StatusForbidden, common.Error{Error: "error: access is not permitted"})
			return
		}

		pass = gen.GenerateUUID()

		data = core.User{
			GroupID:       authResult.Group.ID,
			Name:          input.Name,
			NameEn:        input.NameEn,
			Email:         input.Email,
			Pass:          strings.ToLower(hash.Generate(pass)),
			ExpiredStatus: &[]uint{0}[0],
			Level:         4,
			MailVerify:    &[]bool{false}[0],
			MailToken:     mailToken,
		}
	}

	//check exist for database
	err = dbUser.Create(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	attachment := slack.Attachment{}

	if data.GroupID == 0 {
		attachment.AddField(slack.Field{Title: "Title", Value: "新規ユーザ登録"}).
			AddField(slack.Field{Title: "メールアドレス", Value: input.Email}).
			AddField(slack.Field{Title: "Name", Value: input.Name}).
			AddField(slack.Field{Title: "Name(English)", Value: input.NameEn})
	} else {
		attachment.AddField(slack.Field{Title: "Title", Value: "グループ内ユーザ登録"}).
			AddField(slack.Field{Title: "メールアドレス", Value: input.Email}).
			AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(data.GroupID)) + ":" + authResult.Group.Org}).
			AddField(slack.Field{Title: "Name", Value: input.Name}).
			AddField(slack.Field{Title: "Name(English)", Value: input.NameEn})
	}
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

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

	result := dbUser.Get(user.MailToken, &core.User{MailToken: mailToken})
	if result.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: result.Err.Error() + "| we can't find token data"})
		return
	}

	if *result.User[0].MailVerify {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("This email has already been checked")})
		return
	}
	if *result.User[0].ExpiredStatus >= 1 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: fmt.Sprintf("error: あなたのアカウントは凍結されています。")})
		return
	}

	if err := dbUser.Update(user.UpdateVerifyMail, &core.User{
		Model:      gorm.Model{ID: result.User[0].ID},
		MailVerify: &[]bool{true}[0],
	}); err != nil {
		c.HTML(200, "ng.html", gin.H{})
		//c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.HTML(200, "ok.html", gin.H{})
		//c.JSON(http.StatusOK, &common.Result{Result: "OK"})
	}
}

func Update(c *gin.Context) {
	var input core.User

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("id error")})
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: This GroupID is bad request (0)"})
		return
	}

	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	authResult := auth.UserAuthentication(core.Token{UserToken: userToken, AccessToken: accessToken})
	if authResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: authResult.Err.Error()})
		return
	}

	if !*authResult.User.MailVerify {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "not verify for user mail"})
		return
	}

	var u, serverData core.User

	if authResult.User.ID == uint(id) || id == 0 {
		serverData = authResult.User
	} else {
		if authResult.User.GroupID == 0 {
			c.JSON(http.StatusForbidden, common.Error{Error: "error: Group ID = 0"})
			return
		}
		if authResult.User.Level > 2 {
			c.JSON(http.StatusForbidden, common.Error{Error: "error: failed user level"})
			return
		}
		userResult := dbUser.Get(user.ID, &core.User{Model: gorm.Model{ID: uint(id)}})
		if userResult.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: userResult.Err.Error()})
			return
		}
		if userResult.User[0].GroupID != authResult.User.GroupID {
			c.JSON(http.StatusInternalServerError, common.Error{Error: fmt.Sprintf("failed group authentication")})
			return
		}
		serverData = userResult.User[0]
	}

	u, err = replaceUser(serverData, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if err = dbUser.Update(user.UpdateAll, &u); err != nil {
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

	authResult := auth.UserAuthentication(core.Token{UserToken: userToken, AccessToken: accessToken})
	authResult.User.Pass = ""
	authResult.User.MailToken = ""
	if authResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: authResult.Err.Error()})
		return
	}

	var tmpUser core.User

	if authResult.User.ID == uint(id) {
		tmpUser = authResult.User
	} else if authResult.User.GroupID != 0 {
		if authResult.User.Level >= 2 {
			c.JSON(http.StatusForbidden, common.Error{Error: "You don't have the authority."})
			return
		}

		resultUser := dbUser.Get(user.ID, &core.User{Model: gorm.Model{ID: uint(id)}})
		if resultUser.Err != nil {
			c.JSON(http.StatusUnauthorized, common.Error{Error: resultUser.Err.Error()})
			return
		}

		if resultUser.User[0].GroupID != authResult.User.GroupID {
			c.JSON(http.StatusBadRequest, common.Error{Error: "GroupID is not match."})
			return
		}
		tmpUser = resultUser.User[0]
	}

	c.JSON(http.StatusOK, user.ResultOne{
		ID:         tmpUser.ID,
		GroupID:    tmpUser.GroupID,
		Name:       tmpUser.Name,
		NameEn:     tmpUser.NameEn,
		Email:      tmpUser.Email,
		Level:      tmpUser.Level,
		MailVerify: tmpUser.MailVerify,
	})
}

func GetOwn(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	authResult := auth.UserAuthentication(core.Token{UserToken: userToken, AccessToken: accessToken})
	authResult.User.Pass = ""
	authResult.User.MailToken = ""
	if authResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: authResult.Err.Error()})
	} else {
		c.JSON(http.StatusOK, user.ResultOne{
			ID:         authResult.User.ID,
			GroupID:    authResult.User.GroupID,
			Name:       authResult.User.Name,
			NameEn:     authResult.User.NameEn,
			Email:      authResult.User.Email,
			Level:      authResult.User.Level,
			MailVerify: authResult.User.MailVerify,
		})
	}
}

func GetGroup(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	authUserResult := auth.UserAuthentication(core.Token{UserToken: userToken, AccessToken: accessToken})
	if authUserResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: authUserResult.Err.Error()})
		return
	}

	var data user.Result

	// User権限がLevel=2の時、又はユーザのGroupIDが0の時（グループ未登録時）
	if authUserResult.User.Level > 1 || authUserResult.User.GroupID == 0 {
		data.User = append(data.User, user.ResultOne{
			ID:         authUserResult.User.ID,
			GroupID:    authUserResult.User.GroupID,
			Name:       authUserResult.User.Name,
			NameEn:     authUserResult.User.NameEn,
			Email:      authUserResult.User.Email,
			Level:      authUserResult.User.Level,
			MailVerify: authUserResult.User.MailVerify,
		})
	} else if authUserResult.User.GroupID != 0 {
		resultGroupUser := dbUser.Get(user.GID, &core.User{GroupID: authUserResult.User.GroupID})
		if resultGroupUser.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: resultGroupUser.Err.Error()})
			return
		}

		for _, grp := range resultGroupUser.User {
			data.User = append(data.User, user.ResultOne{
				ID:         grp.ID,
				GroupID:    grp.GroupID,
				Name:       grp.Name,
				NameEn:     grp.NameEn,
				Email:      grp.Email,
				Level:      grp.Level,
				MailVerify: grp.MailVerify,
			})
		}
	}

	c.JSON(http.StatusOK, data)
}
