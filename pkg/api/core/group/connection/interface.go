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
	UpdateAll      = 110
)

type Connection struct {
	gorm.Model
	GroupID   uint   `json:"group_id"`
	ServiceID string `json:"service_id"`
	Service   string `json:"service"`
	NTT       string `json:"ntt"`
	Fee       string `json:"fee"`
	NOC       string `json:"noc"`
	TermIP    string `json:"term_ip"`
	LinkV4    string `json:"link_v4"`
	LinkV6    string `json:"link_v6"`
	Name      string `json:"name"`
	Org       string `json:"org"`
	PostCode  string `json:"postcode"`
	Address   string `json:"address"`
	Mail      string `json:"mail"`
	Phone     string `json:"phone"`
	Country   string `json:"country"`
	Comment   string `json:"comment"`
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
