package ip

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

type Result struct {
	IP []core.IP `json:"ip"`
}

type ResultDatabase struct {
	Err error
	IP  []core.IP
}
