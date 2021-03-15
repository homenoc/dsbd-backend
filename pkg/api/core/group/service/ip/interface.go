package ip

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
	IP []core.IP `json:"ip"`
}

type ResultDatabase struct {
	Err error
	IP  []core.IP
}
