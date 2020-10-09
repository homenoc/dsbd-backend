package auth

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
)

type UserResult struct {
	User user.User
	Err  error
}

type GroupResult struct {
	Group group.Group
	User  user.User
	Err   error
}

type AdminStruct struct {
	User string
	Pass string
}
