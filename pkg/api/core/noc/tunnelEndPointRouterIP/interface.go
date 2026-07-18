package tunnelEndPointRouterIP

import "github.com/homenoc/dsbd-backend/pkg/api/core"

type Result struct {
	TunnelEndPointRouterIP []core.TunnelEndPointRouterIP `json:"gateway_endpoint_ip"`
}

type ResultDatabase struct {
	Err                    error
	TunnelEndPointRouterIP []core.TunnelEndPointRouterIP
}
