package group

import "github.com/jinzhu/gorm"

const (
	ID             = 0
	OrgJa          = 1
	Org            = 2
	Email          = 3
	GID            = 4
	UpdateID       = 100
	UpdateGID      = 101
	UpdateUserInfo = 102
	UpdateTechID   = 103
	UpdateInfo     = 104
	UpdateData     = 105
	UpdateAll      = 110
)

type Connection struct {
	gorm.Model
	GroupID    uint   `json:"group_id"`
	ServiceID  string `json:"service_id"`
	UserId     uint   `json:"user_id"`
	Service    string `json:"service"`
	NTT        string `json:"ntt"`
	NOC        string `json:"noc"`
	NOCIP      string `json:"noc_ip"`
	TermIP     string `json:"term_ip"`
	Monitor    bool   `json:"monitor"`
	LinkV4Our  string `json:"link_v4_our"`
	LinkV4Your string `json:"link_v4_your"`
	LinkV6Our  string `json:"link_v6_our"`
	LinkV6Your string `json:"link_v6_your"`
	Fee        string `json:"fee"`
	Open       bool   `json:"open"`
	Comment    string `json:"comment"`
}

type Result struct {
	Status         bool         `json:"status"`
	Error          string       `json:"error"`
	ConnectionData []Connection `json:"data"`
}

type ResultDatabase struct {
	Err        error
	Connection []Connection
}
