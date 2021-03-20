package jpnicAdmin

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID               = 0
	GroupID          = 1
	UserId           = 2
	NetworkAndUserId = 3
	UpdateLock       = 100
	UpdateAll        = 150
)

type Result struct {
	Admins []core.JPNICAdmin `json:"jpnic_admins"`
}

type ResultDatabase struct {
	Err    error
	Admins []core.JPNICAdmin
}
