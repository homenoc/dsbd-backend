package connection

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
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
