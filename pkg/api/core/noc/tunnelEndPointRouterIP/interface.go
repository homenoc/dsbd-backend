package tunnelEndPointRouterIP

import "github.com/homenoc/dsbd-backend/pkg/api/core"

const (
	ID        = 0
	Enable    = 1
	UpdateAll = 150
)

type Result struct {
	TunnelEndPointRouterIP []core.TunnelEndPointRouterIP `json:"gateway_endpoint_ip"`
}

type ResultDatabase struct {
	Err                    error
	TunnelEndPointRouterIP []core.TunnelEndPointRouterIP
}
