package connection

import (
	"github.com/jinzhu/gorm"
)

const (
	ID              = 0
	OrgJa           = 1
	Org             = 2
	Email           = 3
	GID             = 4
	SearchNewNumber = 5
	UpdateID        = 100
	UpdateGID       = 101
	UpdateUserInfo  = 102
	UpdateTechID    = 103
	UpdateInfo      = 104
	UpdateData      = 105
	UpdateAll       = 110
)

type Connection struct {
	gorm.Model
	NetworkID         uint   `json:"network_id"`
	GroupID           uint   `json:"group_id"`
	UserID            uint   `json:"user_id"`
	RouterID          *uint  `json:"router_id"`     //使用RouterのID
	GatewayIPID       *uint  `json:"gateway_ip_id"` //使用エンドポイントルータのID
	ConnectionType    string `json:"connection_type"`
	ConnectionComment string `json:"connection_comment"` // ServiceがETCの時や補足説明で必要
	ConnectionNumber  uint   `json:"connection_number"`
	NTT               string `json:"ntt"`
	NOC               string `json:"noc"`
	TermIP            string `json:"term_ip"`
	Monitor           *bool  `json:"monitor"`
	LinkV4Our         string `json:"link_v4_our"`
	LinkV4Your        string `json:"link_v4_your"`
	LinkV6Our         string `json:"link_v6_our"`
	LinkV6Your        string `json:"link_v6_your"`
	Fee               string `json:"fee"`
	Open              *bool  `json:"open"`
	Lock              *bool  `json:"lock"`
	Comment           string `json:"comment"`
}

type Result struct {
	Connection []Connection `json:"connection"`
}

type ResultDatabase struct {
	Err        error
	Connection []Connection
}
