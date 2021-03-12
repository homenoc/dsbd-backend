package jpnicTech

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID         = 0
	GroupID    = 1
	UpdateLock = 100
	UpdateAll  = 150
)

type Result struct {
	Tech []core.JPNICTech `json:"tech"`
}

type ResultDatabase struct {
	Err  error
	Tech []core.JPNICTech
}
