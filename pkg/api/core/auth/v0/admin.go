package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbToken "github.com/homenoc/dsbd-backend/pkg/api/store/token/v0"
)

func AdminRadiusAuthorization(data auth.AdminStruct) auth.AdminResult {

	if config.Conf.Controller.Admin.AdminAuth.User == data.User && config.Conf.Controller.Admin.AdminAuth.Pass == data.Pass {
		return auth.AdminResult{AdminID: 0, Err: nil}
	}
	// Todo Radius認証追加予定
	return auth.AdminResult{Err: fmt.Errorf("failed")}
}

func AdminAuthorization(accessToken string) auth.AdminResult {
	tokenResult := dbToken.Get(token.AdminToken, &core.Token{AccessToken: accessToken})
	if tokenResult.Err != nil {
		return auth.AdminResult{Err: tokenResult.Err}
	}
	return auth.AdminResult{Err: nil}
}
