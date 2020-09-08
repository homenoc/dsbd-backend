package store

import (
	"github.com/homenoc/dsbd-backend/pkg/auth"
)

type userResult struct {
	user auth.User
	err  error
}

type allUserResult struct {
	user []auth.User
	err  error
}

type groupResult struct {
	group auth.Group
	err   error
}

type allGroupResult struct {
	group []auth.Group
	err   error
}
