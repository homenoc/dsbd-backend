package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbToken "github.com/homenoc/dsbd-backend/pkg/api/store/token/v0"
)

func AdminRadiusAuthentication(data auth.AdminStruct) auth.AdminResult {

	if config.Conf.Controller.Auth.User == data.User && config.Conf.Controller.Auth.Pass == data.Pass {
		return auth.AdminResult{AdminID: 0, Err: nil}
	}
	// Todo Radius認証追加予定
	return auth.AdminResult{Err: fmt.Errorf("failed")}
}

func AdminAuthentication(accessToken string) auth.AdminResult {
	tokenResult := dbToken.Get(token.AdminToken, &token.Token{AccessToken: accessToken})
	if tokenResult.Err != nil {
		return auth.AdminResult{Err: tokenResult.Err}
	}
	return auth.AdminResult{Err: nil}
}
