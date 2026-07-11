package jpnicAdmin

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

type Result struct {
	Admins []core.JPNICAdmin `json:"jpnic_admins"`
}

type ResultDatabase struct {
	Err    error
	Admins []core.JPNICAdmin
}
