package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	toolToken "github.com/homenoc/dsbd-backend/pkg/tool/token"
	"log"
	"strings"
)

func replaceUser(serverData, input, replace user.User) (user.User, error) {
	updateInfo := 0
	//Name
	if input.Name == "" {
		replace.Name = serverData.Name
	} else {
		replace.Name = input.Name
	}

	//E-Mail
	if input.Email == "" {
		replace.Email = serverData.Email
		replace.MailToken = serverData.MailToken
		replace.MailVerify = serverData.MailVerify
	} else {
		if !strings.Contains(input.Email, "@") {
			return user.User{}, fmt.Errorf("wrong email address")
		}
		tmp := dbUser.Get(user.Email, &user.User{Email: input.Email})
		if tmp.Err != nil {
			return replace, tmp.Err
		}
		if len(tmp.User) != 0 {
			log.Println("error: this email is already registered: " + input.Email)
			return replace, fmt.Errorf("error: this email is already registered")
		}

		mailToken, _ := toolToken.Generate(4)
		replace.Email = input.Email
		replace.MailVerify = false
		replace.MailToken = mailToken
	}

	//Pass
	if input.Pass == "" {
		replace.Pass = serverData.Pass
	} else {
		replace.Pass = input.Pass
	}

	//Org
	if input.Org == "" {
		replace.Org = serverData.Org
	} else {
		replace.Org = input.Org
		updateInfo++
	}

	//PostCode
	if input.PostCode == "" {
		replace.PostCode = serverData.PostCode
	} else {
		replace.PostCode = input.PostCode
		updateInfo++
	}

	//Address
	if input.Address == "" {
		replace.Address = serverData.Address
	} else {
		replace.Address = input.Address
		updateInfo++
	}

	//Phone
	if input.Phone == "" {
		replace.Phone = serverData.Phone
	} else {
		replace.Phone = input.Phone
		updateInfo++
	}

	//Country
	if input.Country == "" {
		replace.Country = serverData.Country
	} else {
		replace.Country = input.Country
		updateInfo++
	}

	//Tech
	if serverData.GID != 0 && serverData.Level <= 1 && input.Status == 1 {
		replace.Tech = input.Tech
	} else {
		replace.Tech = serverData.Tech
	}

	if serverData.Status == 0 && updateInfo == 5 {
		replace.Status = 1
	} else if serverData.Status == 0 && updateInfo < 5 {
		return replace, fmt.Errorf("lack of information")
	}

	return replace, nil
}
