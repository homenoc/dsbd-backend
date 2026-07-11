package jpnicTech

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

type Result struct {
	Tech []core.JPNICTech `json:"tech"`
}

type ResultDatabase struct {
	Err  error
	Tech []core.JPNICTech
}
