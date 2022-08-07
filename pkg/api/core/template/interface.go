package service

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
)

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
	Services                  []config.ServiceTemplate      `json:"services"`
	Connections               []config.ConnectionTemplate   `json:"connections"`
	NTTs                      []string                      `json:"ntts"`
	NOC                       []core.NOC                    `json:"nocs"`
	BGPRouter                 []core.BGPRouter              `json:"bgp_router"`
	TunnelEndPointRouter      []core.TunnelEndPointRouter   `json:"tunnel_endpoint_router"`
	TunnelEndPointRouterIP    []core.TunnelEndPointRouterIP `json:"tunnel_endpoint_router_ip"`
	IPv4                      []string                      `json:"ipv4"`
	IPv6                      []string                      `json:"ipv6"`
	IPv4Route                 []string                      `json:"ipv4_route"`
	IPv6Route                 []string                      `json:"ipv6_route"`
	PaymentMembershipTemplate []config.MembershipTemplate   `json:"payment_membership_template"`
}

type ResultAdmin struct {
	Services                  []config.ServiceTemplate      `json:"services"`
	Connections               []config.ConnectionTemplate   `json:"connections"`
	NTTs                      []string                      `json:"ntts"`
	NOC                       []core.NOC                    `json:"nocs"`
	BGPRouter                 []core.BGPRouter              `json:"bgp_router"`
	TunnelEndPointRouter      []core.TunnelEndPointRouter   `json:"tunnel_endpoint_router"`
	TunnelEndPointRouterIP    []core.TunnelEndPointRouterIP `json:"tunnel_endpoint_router_ip"`
	IPv4                      []string                      `json:"ipv4"`
	IPv6                      []string                      `json:"ipv6"`
	IPv4Route                 []string                      `json:"ipv4_route"`
	IPv6Route                 []string                      `json:"ipv6_route"`
	User                      []core.User                   `json:"user"`
	Group                     []core.Group                  `json:"group"`
	PaymentMembershipTemplate []config.MembershipTemplate   `json:"payment_membership_template"`
	MailTemplate              []config.MailTemplate         `json:"mail_template"`
}
