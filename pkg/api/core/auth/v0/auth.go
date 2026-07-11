package v0

import (
	"errors"
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	dbToken "github.com/homenoc/dsbd-backend/pkg/api/store/token/v0"
	"gorm.io/gorm"
	"log"
	"time"
)

func UserAuthorization(data core.Token) auth.UserResult {
	tokens, err := dbToken.GetSession(data.UserToken, data.AccessToken)
	if len(tokens) == 0 {
		return auth.UserResult{Err: fmt.Errorf("auth failed")}
	}
	if err != nil {
		return auth.UserResult{Err: fmt.Errorf("db error")}
	}

	if *tokens[0].User.ExpiredStatus != core.ExpiredNone {
		return auth.UserResult{Err: fmt.Errorf("deleted this user")}
	}

	go renewProcess(tokens[0])

	return auth.UserResult{User: tokens[0].User, Err: nil}
}

// GroupAuthorization is UserAuthorization plus group-membership checks.
// errorType 0: 未審査の場合はエラーを返す(厳格)　1: 未審査の場合エラーを返さない
func GroupAuthorization(errorType uint, data core.Token) auth.UserResult {
	result := UserAuthorization(data)
	if result.Err != nil {
		return result
	}
	user := result.User

	if user.GroupID == nil {
		return auth.UserResult{Err: fmt.Errorf("no group")}
	}

	// 未審査＋errorType = 0の場合
	if !*user.Group.Pass && errorType == 0 {
		return auth.UserResult{Err: fmt.Errorf("error: unexamined")}
	}
	// アカウント失効時の動作
	if msg := core.ExpiredMessage(*user.Group.ExpiredStatus); msg != "" {
		return auth.UserResult{Err: errors.New(msg)}
	}

	return result
}

func renewProcess(t core.Token) {
	if t.ExpiredAt.UTC().Unix() < time.Now().Add(10*time.Minute).UTC().Unix() {
		result := dbToken.Renew(&core.Token{
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
