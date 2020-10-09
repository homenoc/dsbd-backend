package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/store/group/v0"
	dbToken "github.com/homenoc/dsbd-backend/pkg/store/token/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
)

func UserAuthentication(data token.Token) auth.UserResult {
	resultToken := dbToken.Get(token.UserTokenAndAccessToken, &data)
	if len(resultToken.Token) == 0 {
		return auth.UserResult{Err: fmt.Errorf("auth failed")}
	}
	if resultToken.Err != nil {
		return auth.UserResult{Err: fmt.Errorf("db error")}
	}
	log.Println(resultToken)
	resultUser := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: resultToken.Token[0].UID}})
	if resultUser.Err != nil {
		return auth.UserResult{Err: fmt.Errorf("db error")}
	}
	log.Println(resultUser)
	if 100 <= resultUser.User[0].Status {
		return auth.UserResult{Err: fmt.Errorf("deleted this user")}
	}
	return auth.UserResult{User: resultUser.User[0], Err: nil}
}

func GroupAuthentication(data token.Token) auth.GroupResult {
	resultToken := dbToken.Get(token.UserTokenAndAccessToken, &data)
	if len(resultToken.Token) == 0 {
		return auth.GroupResult{Err: fmt.Errorf("auth failed")}
	}
	if resultToken.Err != nil {
		return auth.GroupResult{Err: fmt.Errorf("db error")}
	}
	resultUser := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: resultToken.Token[0].UID}})
	if resultUser.Err != nil {
		return auth.GroupResult{Err: fmt.Errorf("db error")}
	}
	if resultUser.User[0].Status == 0 || 100 <= resultUser.User[0].Status {
		return auth.GroupResult{Err: fmt.Errorf("db error")}
	}
	resultGroup := dbGroup.Get(group.ID, &group.Group{Model: gorm.Model{ID: resultUser.User[0].GID}})
	if resultGroup.Err != nil {
		return auth.GroupResult{Err: fmt.Errorf("db error")}
	}
	if resultGroup.Group[0].Status < 2 || 1000 <= resultGroup.Group[0].Status {
		return auth.GroupResult{Err: fmt.Errorf("error: group status")}
	}
	return auth.GroupResult{User: resultUser.User[0], Group: resultGroup.Group[0], Err: nil}
}
