package connection

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID              = 0
	OrgJa           = 1
	Org             = 2
	Email           = 3
	ServiceID       = 4
	SearchNewNumber = 5
	NOCID           = 6
	UpdateID        = 100
	UpdateServiceID = 101
	UpdateUserInfo  = 102
	UpdateTechID    = 103
	UpdateInfo      = 104
	UpdateData      = 105
	UpdateAll       = 150
)

type Input struct {
	ConnectionType    string `json:"connection_type"`
	ConnectionComment string `json:"connection_comment"` // ServiceがETCの時や補足説明で必要
	PreferredAP       string `json:"preferred_ap"`
	NTT               string `json:"ntt"`
	Address           string `json:"address"`
	IPv4Route         string `json:"ipv4_route"`
	IPv6Route         string `json:"ipv6_route"`
	NOCID             uint   `json:"noc_id"`
	TermIP            string `json:"term_ip"`
	Monitor           bool   `json:"monitor"`
}

type Connection struct {
	ID                           uint   `json:"id"`
	BGPRouterID                  *uint  `json:"bgp_router_id"`                //使用RouterのID
	BGPRouterName                string `json:"bgp_router_name"`              //使用RouterのID
	TunnelEndPointRouterIPID     *uint  `json:"tunnel_endpoint_router_ip_id"` //使用エンドポイントルータのID
	TunnelEndPointRouterIPIDName string `json:"tunnel_endpoint_router_ip_name"`
	ConnectionTemplateID         *uint  `json:"connection_template_id"`
	ConnectionTemplateName       string `json:"connection_template_name"`
	ConnectionComment            string `json:"connection_comment"` // ServiceがETCの時や補足説明で必要
	ConnectionNumber             uint   `json:"connection_number"`
	NTT                          string `json:"ntt"`
	NOCID                        *uint  `json:"noc_id"`
	NOCName                      string `json:"noc_name"`
	TermIP                       string `json:"term_ip"`
	Monitor                      *bool  `json:"monitor"`
	Address                      string `json:"address"` //都道府県　市町村
	LinkV4Our                    string `json:"link_v4_our"`
	LinkV4Your                   string `json:"link_v4_your"`
	LinkV6Our                    string `json:"link_v6_our"`
	LinkV6Your                   string `json:"link_v6_your"`
	Open                         *bool  `json:"open"`
	Lock                         *bool  `json:"lock"`
}

type Result struct {
	Connection []core.Connection `json:"connection"`
}

type ResultDatabase struct {
	Err        error
	Connection []core.Connection
}
