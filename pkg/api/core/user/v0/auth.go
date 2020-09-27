package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbToken "github.com/homenoc/dsbd-backend/pkg/store/token/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	"github.com/jinzhu/gorm"
)

func authentication(data token.Token) (*user.User, error) {
	resultToken, err := dbToken.Get(token.UserTokenAndAccessToken, &data)
	if err != nil {
		return &user.User{}, fmt.Errorf("db error")
	}
	resultUser, err := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: resultToken.UID}})
	if err != nil {
		return &user.User{}, fmt.Errorf("db error")
	}
	return resultUser, err
}
