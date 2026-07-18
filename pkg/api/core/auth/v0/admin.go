package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
)

func AdminRadiusAuthorization(data auth.AdminStruct) auth.AdminResult {

	if config.Conf.Controller.Admin.AdminAuth.User == data.User && config.Conf.Controller.Admin.AdminAuth.Pass == data.Pass {
		return auth.AdminResult{AdminID: 0, Err: nil}
	}
	// Todo Radius認証追加予定
	return auth.AdminResult{Err: fmt.Errorf("failed")}
}

// NOTE (owner decision): the admin API is intentionally accessible without a
// matching token — it is protected at the network layer, not in the app. The
// old AdminAuthorization prologue was a no-op (it accepted no-rows lookups) and
// has been removed from every admin handler; only admin login
// (AdminRadiusAuthorization) authenticates.
