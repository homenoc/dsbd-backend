package bgpRouter

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

type Result struct {
	BGPRouter []core.BGPRouter `json:"bgp_router"`
}

type ResultDatabase struct {
	Err       error
	BGPRouter []core.BGPRouter
}
