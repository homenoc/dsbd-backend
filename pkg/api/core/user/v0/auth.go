package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbToken "github.com/homenoc/dsbd-backend/pkg/store/token/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
)

func authentication(data token.Token) user.Result {
	resultToken := dbToken.Get(token.UserTokenAndAccessToken, &data)
	if !resultToken.Status {
		return user.Result{Status: false, Error: fmt.Sprintf("db error")}
	}
	resultUser := dbUser.Get(user.ID, &user.User{ID: resultToken.Token[0].UID})
	if !resultToken.Status {
		return user.Result{Status: false, Error: fmt.Sprintf("db error")}
	}
	return resultUser
}
