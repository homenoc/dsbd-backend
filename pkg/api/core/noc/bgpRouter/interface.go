package bgpRouter

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID        = 0
	NOC       = 1
	Address   = 2
	Enable    = 3
	UpdateAll = 150
)

type Result struct {
	BGPRouter []core.BGPRouter `json:"bgp_router"`
}

type ResultDatabase struct {
	Err       error
	BGPRouter []core.BGPRouter
}
