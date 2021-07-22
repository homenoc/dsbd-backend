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
	Services                  []core.ServiceTemplate           `json:"services"`
	Connections               []core.ConnectionTemplate        `json:"connections"`
	NTTs                      []core.NTTTemplate               `json:"ntts"`
	NOC                       []core.NOC                       `json:"nocs"`
	BGPRouter                 []core.BGPRouter                 `json:"bgp_router"`
	TunnelEndPointRouter      []core.TunnelEndPointRouter      `json:"tunnel_endpoint_router"`
	TunnelEndPointRouterIP    []core.TunnelEndPointRouterIP    `json:"tunnel_endpoint_router_ip"`
	IPv4                      []core.IPv4Template              `json:"ipv4"`
	IPv6                      []core.IPv6Template              `json:"ipv6"`
	IPv4Route                 []core.IPv4RouteTemplate         `json:"ipv4_route"`
	IPv6Route                 []core.IPv6RouteTemplate         `json:"ipv6_route"`
	PaymentMembershipTemplate []core.PaymentMembershipTemplate `json:"payment_membership_template"`
	PaymentDonateTemplate     []core.PaymentDonateTemplate     `json:"payment_donate_template"`
	PaymentCouponTemplate     []core.PaymentCouponTemplate     `json:"payment_coupon_template"`
}

type ResultAdmin struct {
	Services                  []core.ServiceTemplate           `json:"services"`
	Connections               []core.ConnectionTemplate        `json:"connections"`
	NTTs                      []core.NTTTemplate               `json:"ntts"`
	NOC                       []core.NOC                       `json:"nocs"`
	BGPRouter                 []core.BGPRouter                 `json:"bgp_router"`
	TunnelEndPointRouter      []core.TunnelEndPointRouter      `json:"tunnel_endpoint_router"`
	TunnelEndPointRouterIP    []core.TunnelEndPointRouterIP    `json:"tunnel_endpoint_router_ip"`
	IPv4                      []core.IPv4Template              `json:"ipv4"`
	IPv6                      []core.IPv6Template              `json:"ipv6"`
	IPv4Route                 []core.IPv4RouteTemplate         `json:"ipv4_route"`
	IPv6Route                 []core.IPv6RouteTemplate         `json:"ipv6_route"`
	User                      []core.User                      `json:"user"`
	Group                     []core.Group                     `json:"group"`
	PaymentMembershipTemplate []core.PaymentMembershipTemplate `json:"payment_membership_template"`
	PaymentDonateTemplate     []core.PaymentDonateTemplate     `json:"payment_donate_template"`
	PaymentCouponTemplate     []core.PaymentCouponTemplate     `json:"payment_coupon_template"`
}

type ResultDatabase struct {
	Err      error
	Services []core.ServiceTemplate
}
