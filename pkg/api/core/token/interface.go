package token

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

type Result struct {
	Token []core.Token `json:"token"`
}

type ResultTmpToken struct {
	Token string `json:"token"`
}

// ResultDatabase is retained only for GetAll; new store functions return
// ([]core.Token, error) directly.
type ResultDatabase struct {
	Err   error
	Token []core.Token
}
