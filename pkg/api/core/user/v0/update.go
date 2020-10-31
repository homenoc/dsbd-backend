package v0

import (
	"fmt"
	toolToken "github.com/homenoc/dsbd-backend/pkg/api/core/tool/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"log"
	"strings"
)

func replaceUser(serverData, input user.User) (user.User, error) {
	updateInfo := 0
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
			return user.User{}, fmt.Errorf("wrong email address")
		}
		tmp := dbUser.Get(user.Email, &user.User{Email: input.Email})
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
	}

	//Pass
	if input.Pass != "" {
		serverData.Pass = input.Pass
	}

	//GroupHandle
	if input.GroupHandle != serverData.GroupHandle {
		serverData.GroupHandle = input.GroupHandle
	}

	//Org
	if input.Org != "" {
		serverData.Org = input.Org
		updateInfo++
	}

	//Org (English)
	if input.OrgEn != "" {
		serverData.OrgEn = input.OrgEn
		updateInfo++
	}

	//PostCode
	if input.PostCode != "" {
		serverData.PostCode = input.PostCode
		updateInfo++
	}

	//Address
	if input.Address != "" {
		serverData.Address = input.Address
		updateInfo++
	}

	//Address(English)
	if input.AddressEn != "" {
		serverData.AddressEn = input.AddressEn
		updateInfo++
	}

	//Dept
	if input.Dept != "" {
		serverData.Dept = input.Dept
		updateInfo++
	}

	//Dept(English)
	if input.DeptEn != "" {
		serverData.DeptEn = input.DeptEn
		updateInfo++
	}

	//Pos
	if input.Pos != "" {
		serverData.Pos = input.Pos
		updateInfo++
	}

	//Pos(English)
	if input.PosEn != "" {
		serverData.PosEn = input.PosEn
		updateInfo++
	}

	//Tel
	if input.Tel != "" {
		serverData.Tel = input.Tel
		updateInfo++
	}

	//Fax
	if input.Fax != "" {
		serverData.Fax = input.Fax
		updateInfo++
	}

	//Country
	if input.Country != "" {
		serverData.Country = input.Country
		updateInfo++
	}

	//Tech
	if serverData.GroupID != 0 && serverData.Level <= 1 && input.Status == 1 {
		serverData.Tech = input.Tech
	}

	if serverData.Status == 0 && updateInfo == 12 {
		serverData.Status = 1
	}

	return serverData, nil
}

func updateAdminUser(input, replace user.User) (user.User, error) {
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
		replace.MailVerify = &[]bool{false}[0]
		replace.MailToken = mailToken
	}

	//Pass
	if input.Pass != "" {
		replace.Pass = input.Pass
	}

	//Org
	if input.Org != "" {
		replace.Org = input.Org
	}

	//Org (English)
	if input.OrgEn != "" {
		replace.OrgEn = input.OrgEn
	}

	//PostCode
	if input.PostCode != "" {
		replace.PostCode = input.PostCode
	}

	//Address
	if input.Address != "" {
		replace.Address = input.Address
	}

	//Address(English)
	if input.AddressEn != "" {
		replace.AddressEn = input.AddressEn
	}

	//Dept
	if input.Dept != "" {
		replace.Dept = input.Dept
	}

	//Dept(English)
	if input.DeptEn != "" {
		replace.DeptEn = input.DeptEn
	}

	//Pos
	if input.Pos != "" {
		replace.Pos = input.Pos
	}

	//Pos(English)
	if input.PosEn != "" {
		replace.PosEn = input.PosEn
	}

	//Tel
	if input.Tel != "" {
		replace.Tel = input.Tel
	}

	//Fax
	if input.Fax != "" {
		replace.Fax = input.Fax
	}

	//Country
	if input.Country != "" {
		replace.Country = input.Country
	}

	// uint boolean
	//Tech
	if input.Tech != replace.Tech {
		replace.Tech = input.Tech
	}

	//Status
	if input.Status != replace.Status {
		replace.Status = input.Status
	}

	//Level
	if input.Level != replace.Level {
		replace.Level = input.Level
	}

	return replace, nil
}
