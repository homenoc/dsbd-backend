package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	toolToken "github.com/homenoc/dsbd-backend/pkg/tool/token"
	"log"
	"strconv"
	"strings"
)

func replaceUser(serverData, input, replace user.User) (user.User, error) {
	log.Println("input")
	log.Println(input)

	updateInfo := 0
	//Name
	if input.Name == "" {
		replace.Name = serverData.Name
	} else {
		replace.Name = input.Name
	}

	log.Println(1)

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

	log.Println(2)

	//Pass
	if input.Pass == "" {
		replace.Pass = serverData.Pass
	} else {
		replace.Pass = input.Pass
	}
	log.Println(3)

	//Org
	if input.Org == "" {
		replace.Org = serverData.Org
	} else {
		replace.Org = input.Org
		updateInfo++
	}
	log.Println(4)

	//PostCode
	if input.PostCode == "" {
		replace.PostCode = serverData.PostCode
	} else {
		replace.PostCode = input.PostCode
		updateInfo++
	}

	log.Println(5)

	//Address
	if input.Address == "" {
		replace.Address = serverData.Address
	} else {
		replace.Address = input.Address
		updateInfo++
	}

	log.Println(6)

	//Phone
	if input.Phone == "" {
		replace.Phone = serverData.Phone
	} else {
		replace.Phone = input.Phone
		updateInfo++
	}
	log.Println(7)

	//Country
	if input.Country == "" {
		replace.Country = serverData.Country
	} else {
		replace.Country = input.Country
		updateInfo++
	}
	log.Println(10)
	log.Println("updateinfo: " + strconv.Itoa(updateInfo))

	if serverData.Status == 0 && updateInfo == 5 {
		replace.Status = 1
	} else if serverData.Status == 0 && updateInfo < 5 {
		return replace, fmt.Errorf("lack of information")
	}

	return replace, nil
}
