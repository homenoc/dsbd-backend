package catalog

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
)

// Result is the frontend bootstrap catalog: service/connection type definitions
// (from the code registry in package core) plus the deployment option lists
// (from config). It replaces the non-entity part of the old /template blob;
// entity lists (NOC/routers/users/groups) are fetched from their own endpoints.
type Result struct {
	Services          []core.ServiceType          `json:"services"`
	Connections       []core.ConnectionType       `json:"connections"`
	IX                []config.IXTemplate         `json:"ix"`
	NTTs              []string                    `json:"ntts"`
	IPv4              []string                    `json:"ipv4"`
	IPv6              []string                    `json:"ipv6"`
	IPv4Route         []string                    `json:"ipv4_route"`
	IPv6Route         []string                    `json:"ipv6_route"`
	PreferredAP       []string                    `json:"preferred_ap"`
	PaymentMembership []config.MembershipTemplate `json:"payment_membership"`
	MailTemplate      []config.MailTemplate       `json:"mail_template"`
	MemberType        []core.ConstantMembership   `json:"member_type"`
}
