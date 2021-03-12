package tunnelEndPointRouterIP

import "github.com/homenoc/dsbd-backend/pkg/api/core"

const (
	ID        = 0
	NOC       = 1
	Address   = 2
	Enable    = 3
	UpdateAll = 150
)

type Result struct {
	TunnelEndPointRouterIP []core.TunnelEndPointRouterIP `json:"gateway_endpoint_ip"`
}

type ResultDatabase struct {
	Err                    error
	TunnelEndPointRouterIP []core.TunnelEndPointRouterIP
}
