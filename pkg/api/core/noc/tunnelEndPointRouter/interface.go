package tunnelEndPointRouter

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
	TunnelEndPointRouters []core.TunnelEndPointRouter `json:"tunnel_endpoint_routers"`
}

type ResultDatabase struct {
	Err                  error
	TunnelEndPointRouter []core.TunnelEndPointRouter
}
