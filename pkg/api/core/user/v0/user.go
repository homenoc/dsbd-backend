package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/mail"
	"github.com/homenoc/dsbd-backend/pkg/api/core/mail/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/gen"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/hash"
	toolToken "github.com/homenoc/dsbd-backend/pkg/api/core/tool/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Add(c *gin.Context) {
	var input user.Input
	var data core.User

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

	// 新規ユーザ
	if input.Pass == "" {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("wrong pass")})
		return
	}

	data = core.User{
		GroupID:       nil,
		Name:          input.Name,
		NameEn:        input.NameEn,
		Email:         input.Email,
		Pass:          input.Pass,
		ExpiredStatus: &[]uint{core.ExpiredNone}[0],
		Level:         core.LevelMaster,
		MailVerify:    &[]bool{false}[0],
		MailToken:     mailToken,
	}

	//check exist for database
	result := dbUser.GetByEmail(input.Email)
	if result.Err != nil {
		log.Println(result.Err)
	}

	if len(result.User) != 0 && result.Err == nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: "this email is already registered: \" + u.Email"})
		return
	}

	err = dbUser.Create(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	noticeAdd(input)

	mailTemplate, _ := config.GetMailTemplate("signature")

	v0.SendMail(mail.Mail{
		ToMail:  data.Email,
		Subject: "本人確認のメールにつきまして",
		Content: " " + input.Name + "様\n\n" + "以下のリンクから本人確認を完了してください。\n" +
			config.Conf.Controller.User.Url + "/api/v1/verify/" + mailToken + "\n" +
			"本人確認が完了次第、ログイン可能になります。" + mailTemplate.Message,
	})

	c.JSON(http.StatusOK, user.Result{})
}

func AddGroup(c *gin.Context) {
	var input user.Input
	var data core.User

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("id error")})
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: This GroupID is bad request (0)"})
		return
	}

	err = c.BindJSON(&input)
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

	// グループ所属ユーザの登録
	currentUser := middleware.CurrentUser(c)

	if !core.CanManageServices(currentUser.Level) {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: access is not permitted"})
		return
	}

	resultGroup := dbGroup.GetByID(uint(id))
	if resultGroup.Err != nil {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: access is not permitted"})
		return
	}

	if currentUser.GroupID == nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: group is not found"})
		return
	}

	if *currentUser.GroupID != uint(id) {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: group id is invalid"})
		return
	}

	if !(core.LevelEditor <= input.Level && input.Level <= core.LevelGuest) {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: user level is invalid"})
		return
	}

	pass, err = gen.GenerateUUIDString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "error: Failed to generate uuid. "})
		return
	}

	data = core.User{
		GroupID:       currentUser.GroupID,
		Name:          input.Name,
		NameEn:        input.NameEn,
		Email:         input.Email,
		Pass:          strings.ToLower(hash.Generate(pass)),
		ExpiredStatus: &[]uint{core.ExpiredNone}[0],
		Level:         input.Level,
		MailVerify:    &[]bool{false}[0],
		MailToken:     mailToken,
	}

	//check exist for database
	result := dbUser.GetByEmail(input.Email)
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	if len(result.User) != 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: "this email is already registered: \" + u.Email"})
		return
	}

	err = dbUser.Create(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	noticeAddFromGroup(input, *currentUser.Group)

	mailTemplate, _ := config.GetMailTemplate("signature")

	v0.SendMail(mail.Mail{
		ToMail:  data.Email,
		Subject: "本人確認メールにつきまして",
		Content: " " + input.Name + "様\n\n" + "以下のリンクから本人確認を完了してください。\n" +
			config.Conf.Controller.User.Url + "/api/v1/verify/" + mailToken + "\n" +
			"本人確認が完了次第、ログイン可能になります。\n" + "仮パスワード: " + pass + mailTemplate.Message,
	})

	c.JSON(http.StatusOK, user.Result{})
}

func MailVerify(c *gin.Context) {
	mailToken := c.Param("token")

	result := dbUser.GetByMailToken(mailToken)
	if result.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: result.Err.Error() + "| we can't find token data"})
		return
	}

	if *result.User[0].MailVerify {
		c.Writer.WriteString(`<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <title>dsysメール確認システム</title>
    <meta http-equiv="refresh" content="5; URL=` + config.Conf.Web.URL + `">
</head>
<body>
<h1>すでにメールアドレスの確認はできています。</h1>
<p>5秒後にログイン画面に移動します</p>
<br>
<p>This email has already been checked</p>
</body>
</html>`)
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("This email has already been checked")})
		return
	}
	if *result.User[0].ExpiredStatus >= 1 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: fmt.Sprintf("error: あなたのアカウントは凍結されています。")})
		return
	}

	if err := dbUser.UpdateMailVerify(&core.User{
		Model:      gorm.Model{ID: result.User[0].ID},
		MailVerify: &[]bool{true}[0],
	}); err != nil {
		c.Writer.WriteString(`<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <title>dsysメール確認システム</title>
    <meta http-equiv="refresh" content="30; URL=` + config.Conf.Web.URL + `">
</head>
<body>
<h1>メールの確認ができませんでした。</h1>
<br>
<p>` + err.Error() +
			`</p>
</body>
</html>`)
	} else {
		c.Writer.WriteString(`<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <title>dsysメール確認システム</title>
    <meta http-equiv="refresh" content="5; URL=` + config.Conf.Web.URL + `">
</head>
<body>
<h1>メールの確認ができました。</h1>
<p>5秒後にログイン画面に移動します</p>
</body>
</html>`)
	}
}

func Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("id error")})
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: This GroupID is bad request (0)"})
		return
	}

	currentUser := middleware.CurrentUser(c)

	if !core.CanViewGroup(currentUser.Level) {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: failed user level"})
		return
	}

	u := dbUser.GetByID(uint(id))
	if u.Err != nil {
		log.Println(u.Err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: "error: database error"})
		return
	}

	if u.User[0].GroupID == nil || *u.User[0].GroupID != *currentUser.GroupID {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: This user does not belong to your group."})
		return
	}

	if u.User[0].Level < core.LevelEditor {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: The master user cannot be deleted."})
		return
	}

	err = dbUser.Delete(&core.User{Model: gorm.Model{ID: uint(id)}})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: delete error. "})
		return
	}

	c.JSON(http.StatusOK, common.Result{})
}

func Update(c *gin.Context) {
	var input user.Input

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("id error")})
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: This GroupID is bad request (0)"})
		return
	}

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	currentUser := middleware.CurrentUser(c)

	var u, serverData core.User

	if currentUser.ID == uint(id) || id == 0 {
		serverData = currentUser
	} else {
		if currentUser.GroupID == nil {
			c.JSON(http.StatusForbidden, common.Error{Error: "error: Group ID = 0"})
			return
		}
		// Master/Member（Level 1,2）のみユーザ設定を変更可能。Level 3以上は不可
		if !core.CanManageServices(currentUser.Level) {
			c.JSON(http.StatusForbidden, common.Error{Error: "error: failed user level"})
			return
		}
		userResult := dbUser.GetByID(uint(id))
		if userResult.Err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: userResult.Err.Error()})
			return
		}

		if userResult.User[0].GroupID == nil || *userResult.User[0].GroupID != *currentUser.GroupID {
			c.JSON(http.StatusBadRequest, common.Error{Error: "error: This user does not belong to your group."})
			return
		}

		serverData = userResult.User[0]
	}

	u, err = replaceUser(serverData, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	noticeRenew(currentUser, serverData, input)

	if err = dbUser.Update(&u); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, user.Result{})
	}
}
