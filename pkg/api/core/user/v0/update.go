package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/mail"
	"github.com/homenoc/dsbd-backend/pkg/api/core/mail/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	toolToken "github.com/homenoc/dsbd-backend/pkg/api/core/tool/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"log"
	"strings"
)

func replaceUser(serverData core.User, input user.Input) (core.User, error) {
	//Name
	if input.Name != "" {
		serverData.Name = input.Name
	}

	//Name (English)
	if input.NameEn != "" {
		serverData.NameEn = input.NameEn
	}

	//E-Mail
	if input.Email != "" {
		if !strings.Contains(input.Email, "@") {
			return core.User{}, fmt.Errorf("wrong email address")
		}
		tmp := dbUser.Get(user.Email, &core.User{Email: input.Email})
		if tmp.Err != nil {
			return serverData, tmp.Err
		}
		if len(tmp.User) != 0 {
			log.Println("error: this email is already registered: " + input.Email)
			return serverData, fmt.Errorf("error: this email is already registered")
		}

		mailToken, _ := toolToken.Generate(4)
		serverData.Email = input.Email
		serverData.MailVerify = &[]bool{false}[0]
		serverData.MailToken = mailToken

		v0.SendMail(mail.Mail{
			ToMail:  input.Email,
			Subject: "本人確認のメールにつきまして",
			Content: " " + serverData.Name + "様\n\n" + "以下のリンクから本人確認を完了してください。\n" +
				config.Conf.Controller.User.Url + "/api/v1/verify/" + mailToken + "\n" +
				"本人確認が完了次第、ログイン可能になります。\n",
		})
	}

	//Pass
	if input.Pass != "" {
		serverData.Pass = input.Pass

		v0.SendMail(mail.Mail{
			ToMail:  serverData.Email,
			Subject: "[通知] パスワード変更",
			Content: " " + serverData.Name + "様\n\n" + "パスワードが変更されました。\n",
		})
	}

	//Level
	if input.Level != 0 {
		if !(1 < input.Level && input.Level < 5) {
			return core.User{}, fmt.Errorf("error: user level is invalid")
		} else {
			serverData.Level = input.Level
		}
	}

	return serverData, nil
}

func updateAdminUser(input, replace core.User) (core.User, error) {
	//Name
	if input.Name != "" {
		replace.Name = input.Name
	}

	//Name (English)
	if input.NameEn != "" {
		replace.NameEn = input.NameEn
	}

	//E-Mail
	if input.Email != "" {
		if !strings.Contains(input.Email, "@") {
			return core.User{}, fmt.Errorf("wrong email address")
		}
		tmp := dbUser.Get(user.Email, &core.User{Email: input.Email})
		if tmp.Err != nil {
			return replace, tmp.Err
		}
		if len(tmp.User) != 0 {
			log.Println("error: this email is already registered: " + input.Email)
			return replace, fmt.Errorf("error: this email is already registered")
		}

		mailToken, _ := toolToken.Generate(4)
		replace.Email = input.Email
		replace.MailVerify = &[]bool{false}[0]
		replace.MailToken = mailToken
	}

	//Pass
	if input.Pass != "" {
		replace.Pass = input.Pass
	}

	// uint boolean

	//Level
	if input.Level != replace.Level {
		replace.Level = input.Level
	}

	//GroupID
	if input.GroupID != replace.GroupID {
		replace.GroupID = input.GroupID
	}

	return replace, nil
}
