package service

import "github.com/homenoc/dsbd-backend/pkg/api/core"

const (
	ID        = 0
	UpdateAll = 150
)

type Input struct {
	Hidden  *bool  `json:"hidden"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type Result struct {
	Services               []core.ServiceTemplate        `json:"services"`
	Connections            []core.ConnectionTemplate     `json:"connections"`
	NTTs                   []core.NTTTemplate            `json:"ntts"`
	NOC                    []core.NOC                    `json:"nocs"`
	BGPRouter              []core.BGPRouter              `json:"bgp_router"`
	TunnelEndPointRouter   []core.TunnelEndPointRouter   `json:"tunnel_endpoint_router"`
	TunnelEndPointRouterIP []core.TunnelEndPointRouterIP `json:"tunnel_endpoint_router_ip"`
}

type ResultDatabase struct {
	Err      error
	Services []core.ServiceTemplate
}
