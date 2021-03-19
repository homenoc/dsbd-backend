package auth

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

type UserResult struct {
	User core.User
	Err  error
}

type GroupResult struct {
	User core.User
	Err  error
}

type AdminStruct struct {
	User string
	Pass string
}

type AdminResult struct {
	AdminID uint
	Err     error
}
