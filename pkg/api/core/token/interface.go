package token

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID                      = 0
	UserToken               = 10
	UserTokenAndAccessToken = 11
	ExpiredTime             = 12
	AdminToken              = 20
	AddToken                = 100
	UpdateToken             = 101
	UpdateAll               = 150
)

type Result struct {
	Token []core.Token `json:"token"`
}

type ResultTmpToken struct {
	Token string `json:"token"`
}

type ResultDatabase struct {
	Err   error
	Token []core.Token
}
