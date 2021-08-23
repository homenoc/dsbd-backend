package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbToken "github.com/homenoc/dsbd-backend/pkg/api/store/token/v0"
	"gorm.io/gorm"
	"log"
	"time"
)

func UserAuthentication(data core.Token) auth.UserResult {
	resultToken := dbToken.Get(token.UserTokenAndAccessToken, &data)
	if len(resultToken.Token) == 0 {
		return auth.UserResult{Err: fmt.Errorf("auth failed")}
	}
	if resultToken.Err != nil {
		return auth.UserResult{Err: fmt.Errorf("db error")}
	}

	if 0 < *resultToken.Token[0].User.ExpiredStatus {
		return auth.UserResult{Err: fmt.Errorf("deleted this user")}
	}

	go renewProcess(resultToken.Token[0])

	return auth.UserResult{User: resultToken.Token[0].User, Err: nil}
}

// errorType 0: 未審査の場合はエラーを返す(厳格)　1: 未審査の場合エラーを返さない
func GroupAuthentication(errorType uint, data core.Token) auth.GroupResult {
	resultToken := dbToken.Get(token.UserTokenAndAccessToken, &data)
	if len(resultToken.Token) == 0 {
		return auth.GroupResult{Err: fmt.Errorf("auth failed")}
	}
	if resultToken.Err != nil {
		return auth.GroupResult{Err: fmt.Errorf("error: no token")}
	}

	if 0 < *resultToken.Token[0].User.ExpiredStatus {
		return auth.GroupResult{Err: fmt.Errorf("deleted this user")}
	}

	if resultToken.Token[0].User.GroupID == nil {
		return auth.GroupResult{Err: fmt.Errorf("no group")}
	}

	// 未審査＋errorType = 0の場合
	if !*resultToken.Token[0].User.Group.Pass && errorType == 0 {
		return auth.GroupResult{Err: fmt.Errorf("error: unexamined")}
	}
	// アカウント失効時の動作
	if *resultToken.Token[0].User.Group.ExpiredStatus == 1 {
		return auth.GroupResult{Err: fmt.Errorf("error: discontinued by Master Account")}
	}
	if *resultToken.Token[0].User.Group.ExpiredStatus == 2 {
		return auth.GroupResult{Err: fmt.Errorf("error: discontinuation by the steering committee")}
	}
	if *resultToken.Token[0].User.Group.ExpiredStatus == 3 {
		return auth.GroupResult{Err: fmt.Errorf("error: discontinuation due to failed review")}
	}

	go renewProcess(resultToken.Token[0])

	return auth.GroupResult{User: resultToken.Token[0].User, Err: nil}
}

func renewProcess(t core.Token) {
	if t.ExpiredAt.UTC().Unix() < time.Now().Add(10*time.Minute).UTC().Unix() {
		result := dbToken.Update(token.UpdateToken, &core.Token{
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
