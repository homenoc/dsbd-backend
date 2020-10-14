package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	toolToken "github.com/homenoc/dsbd-backend/pkg/tool/token"
	"strings"
)

func replaceUser(authResult, input, replace user.User) (user.User, error) {
	//Name
	if input.Name == "" {
		replace.Name = authResult.Name
	} else {
		replace.Name = input.Name
	}

	//E-Mail
	if !strings.Contains(input.Email, "@") {
		return user.User{}, fmt.Errorf("wrong email address")
	}
	if input.Email == "" {
		replace.Email = authResult.Email
		replace.MailToken = authResult.MailToken
		replace.MailVerify = authResult.MailVerify
	} else {
		mailToken, _ := toolToken.Generate(4)
		replace.Email = input.Email
		replace.MailVerify = false
		replace.MailToken = mailToken
	}

	//Pass
	if input.Pass == "" {
		replace.Pass = authResult.Pass
	} else {
		replace.Pass = input.Pass
	}

	//Org
	if input.Org == "" {
		replace.Org = authResult.Org
	} else {
		replace.Org = input.Org
	}

	//PostCode
	if input.PostCode == "" {
		replace.PostCode = authResult.PostCode
	} else {
		replace.PostCode = input.PostCode
	}

	//Address
	if input.Address == "" {
		replace.Address = authResult.Address
	} else {
		replace.Address = input.Address
	}

	//Phone
	if input.Phone == "" {
		replace.Phone = authResult.Phone
	} else {
		replace.Phone = input.Phone
	}

	//Country
	if input.Country == "" {
		replace.Country = authResult.Country
	} else {
		replace.Country = input.Country
	}

	return replace, nil
}
