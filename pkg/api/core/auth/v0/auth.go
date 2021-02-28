package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbToken "github.com/homenoc/dsbd-backend/pkg/api/store/token/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func UserAuthentication(data token.Token) auth.UserResult {
	resultToken := dbToken.Get(token.UserTokenAndAccessToken, &data)
	if len(resultToken.Token) == 0 {
		return auth.UserResult{Err: fmt.Errorf("auth failed")}
	}
	if resultToken.Err != nil {
		return auth.UserResult{Err: fmt.Errorf("db error")}
	}
	resultUser := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: resultToken.Token[0].UserID}})
	if resultUser.Err != nil {
		return auth.UserResult{Err: fmt.Errorf("db error")}
	}
	if 100 <= resultUser.User[0].Status {
		return auth.UserResult{Err: fmt.Errorf("deleted this user")}
	}

	renewProcess(resultToken.Token[0])

	return auth.UserResult{User: resultUser.User[0], Err: nil}
}

// errorType 0: 未審査の場合でもエラーを返す　1: 未審査の場合エラーを返さない
func GroupAuthentication(errorType uint, data token.Token) auth.GroupResult {
	resultToken := dbToken.Get(token.UserTokenAndAccessToken, &data)
	if len(resultToken.Token) == 0 {
		return auth.GroupResult{Err: fmt.Errorf("auth failed")}
	}
	if resultToken.Err != nil {
		return auth.GroupResult{Err: fmt.Errorf("error: no token")}
	}
	resultUser := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: resultToken.Token[0].UserID}})
	if resultUser.Err != nil {
		return auth.GroupResult{Err: fmt.Errorf("db error")}
	}
	if resultUser.User[0].Status == 0 || 100 <= resultUser.User[0].Status {
		return auth.GroupResult{Err: fmt.Errorf("user status error")}
	}
	if resultUser.User[0].GroupID == 0 {
		return auth.GroupResult{Err: fmt.Errorf("no group")}
	}
	resultGroup := dbGroup.Get(group.ID, &group.Group{Model: gorm.Model{ID: resultUser.User[0].GroupID}})
	if resultGroup.Err != nil {
		return auth.GroupResult{Err: fmt.Errorf("db error")}
	}
	// 未審査＋errorType = 0の場合
	if !*resultGroup.Group[0].Pass && errorType == 0 {
		return auth.GroupResult{Err: fmt.Errorf("error: unexamined")}
	}
	// アカウント失効時の動作
	if *resultGroup.Group[0].ExpiredStatus == 1 {
		return auth.GroupResult{Err: fmt.Errorf("error: discontinued by Master Account")}
	}
	if *resultGroup.Group[0].ExpiredStatus == 2 {
		return auth.GroupResult{Err: fmt.Errorf("error: discontinuation by the steering committee")}
	}
	if *resultGroup.Group[0].ExpiredStatus == 3 {
		return auth.GroupResult{Err: fmt.Errorf("error: discontinuation due to failed review")}
	}

	renewProcess(resultToken.Token[0])

	return auth.GroupResult{User: resultUser.User[0], Group: resultGroup.Group[0], Err: nil}
}

func renewProcess(t token.Token) {
	if t.ExpiredAt.UTC().Unix() < time.Now().Add(10*time.Minute).UTC().Unix() {
		result := dbToken.Update(token.UpdateToken, &token.Token{
			Model:     gorm.Model{ID: t.ID},
			ExpiredAt: t.ExpiredAt.Add(10 * time.Minute),
		})
		if err := result; err != nil {
			log.Println(err)
		} else {
			log.Println("Success!!")
		}
	}
}
