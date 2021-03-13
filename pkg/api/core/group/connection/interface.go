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
	UpdateID        = 100
	UpdateServiceID = 101
	UpdateUserInfo  = 102
	UpdateTechID    = 103
	UpdateInfo      = 104
	UpdateData      = 105
	UpdateAll       = 150
)

type Input struct {
	TunnelRouterIPID     *uint  `json:"tunnel_router_ip_id"` //使用エンドポイントルータのID
	ConnectionTemplateID *uint  `json:"connection_template_id"`
	ConnectionComment    string `json:"connection_comment"` // ServiceがETCの時や補足説明で必要
	NTTTemplateID        *uint  `json:"ntt_template_id"`
	NOCID                *uint  `json:"noc_id"`
	TermIP               string `json:"term_ip"`
	Monitor              *bool  `json:"monitor"`
	Prefectures          uint   `json:"prefectures"` //JIS X 0401
}

type Result struct {
	Connection []core.Connection `json:"connection"`
}

type ResultDatabase struct {
	Err        error
	Connection []core.Connection
}
