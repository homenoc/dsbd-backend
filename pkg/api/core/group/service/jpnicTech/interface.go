package jpnicTech

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID               = 0
	NetworkID        = 1
	UserId           = 2
	NetworkAndUserId = 3
	UpdateLock       = 100
	UpdateAll        = 150
)

type Result struct {
	Tech []core.JPNICTech `json:"tech"`
}

type ResultDatabase struct {
	Err  error
	Tech []core.JPNICTech
}
