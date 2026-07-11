package connection

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
)

type Input struct {
	ConnectionType    string `json:"connection_type"`
	ConnectionComment string `json:"connection_comment"` // ServiceがETCの時や補足説明で必要
	PreferredAP       string `json:"preferred_ap"`
	NTT               string `json:"ntt"`
	IX                string `json:"ix"`           // 接続IX
	IXPeerType        string `json:"ix_peer_type"` // パブリック or PI/CUG
	IXVlanID          string `json:"ix_vlan_id"`   // VLAN-ID（PI/CUGの場合）
	LinkV4Your        string `json:"link_v4_your"` // 相手側IPv4アドレス（IX接続時）
	LinkV6Your        string `json:"link_v6_your"` // 相手側IPv6アドレス（IX接続時）
	Address           string `json:"address"`
	IPv4Route         string `json:"ipv4_route"`
	IPv6Route         string `json:"ipv6_route"`
	TermIP            string `json:"term_ip"`
	RFC8950           bool   `json:"rfc8950"`
	Monitor           bool   `json:"monitor"`
	Comment           string `json:"comment"`
}

type Result struct {
	Connection []core.Connection `json:"connection"`
}

type ResultDatabase struct {
	Err        error
	Connection []core.Connection
}

// UserView is the user-facing wire shape of an enabled connection
// (GET /connection).
type UserView struct {
	ID        uint   `json:"id"`
	ServiceID string `json:"service_id"`
	Open      bool   `json:"open"`
}

// NetworkInfo is the derived network summary of an opened connection
// (GET /info) — the contract-disclosure data the web Info page renders.
type NetworkInfo struct {
	ServiceID      string                  `json:"service_id"`
	Service        string                  `json:"service"`
	Assign         bool                    `json:"assign"`
	ASN            uint                    `json:"asn"`
	V4             []string                `json:"v4"`
	V6             []string                `json:"v6"`
	NOC            string                  `json:"noc"`
	NOCIP          string                  `json:"noc_ip"`
	TermIP         string                  `json:"term_ip"`
	RFC8950        bool                    `json:"rfc8950"`
	LinkV4Our      string                  `json:"link_v4_our"`
	LinkV4Your     string                  `json:"link_v4_your"`
	LinkV6Our      string                  `json:"link_v6_our"`
	LinkV6Your     string                  `json:"link_v6_your"`
	Fee            string                  `json:"fee"`
	Org            string                  `json:"org"`
	OrgEn          string                  `json:"org_en"`
	PostCode       string                  `json:"postcode"`
	Address        string                  `json:"address"`
	AddressEn      string                  `json:"address_en"`
	JPNICAdmin     service.UserJPNICView   `json:"jpnic_admin"`
	JPNICTech      []service.UserJPNICView `json:"jpnic_tech"`
	AveUpstream    uint                    `json:"avg_upstream"`
	MaxUpstream    uint                    `json:"max_upstream"`
	AveDownstream  uint                    `json:"avg_downstream"`
	MaxDownstream  uint                    `json:"max_downstream"`
	MaxBandWidthAS string                  `json:"max_bandwidth_as"`
	BGPRouteV4     string                  `json:"bgp_route_v4"`
	BGPRouteV6     string                  `json:"bgp_route_v6"`
}
