package tunnelEndPointRouter

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID        = 0
	UpdateAll = 150
)

type Result struct {
	TunnelEndPointRouters []core.TunnelEndPointRouter `json:"tunnel_endpoint_routers"`
}

type ResultDatabase struct {
	Err                  error
	TunnelEndPointRouter []core.TunnelEndPointRouter
}
