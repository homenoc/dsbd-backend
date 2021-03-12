package notice

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID               = 0
	UserID           = 1
	GroupID          = 2
	UserIDAndGroupID = 3
	Everyone         = 4
	GroupData        = 5
	UserData         = 6
	Important        = 10
	Fault            = 11
	Info             = 12
	UpdateAll        = 150
)

type Result struct {
	Notice []core.Notice `json:"notice"`
}

type ResultDatabase struct {
	Err    error
	Notice []core.Notice
}
