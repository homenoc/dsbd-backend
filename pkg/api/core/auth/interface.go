package auth

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

// UserResult is returned by both UserAuthorization and GroupAuthorization
// (the latter adds group checks on top of the same shape).
type UserResult struct {
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
