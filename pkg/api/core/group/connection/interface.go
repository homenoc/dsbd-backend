package connection

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID              = 0
	OrgJa           = 1
	Org             = 2
	Email           = 3
	ServiceID       = 4
	SearchNewNumber = 5
	UpdateID        = 100
	UpdateServiceID = 101
	UpdateUserInfo  = 102
	UpdateTechID    = 103
	UpdateInfo      = 104
	UpdateData      = 105
	UpdateAll       = 150
)

type Result struct {
	Connection []core.Connection `json:"connection"`
}

type ResultDatabase struct {
	Err        error
	Connection []core.Connection
}
