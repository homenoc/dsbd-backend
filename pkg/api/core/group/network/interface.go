package group

import "github.com/jinzhu/gorm"

const (
	ID          = 0
	GID         = 1
	Org         = 2
	Type        = 3
	UpdateName  = 100
	UpdateDate  = 102
	UpdateRoute = 103
	UpdatePlan  = 104
	UpdateGID   = 104
	UpdateData  = 105
	UpdateAll   = 110
)

type Network struct {
	gorm.Model
	GroupID   uint   `json:"group_id"`
	Org       string `json:"org"`
	OrgEn     string `json:"org_en"`
	Postcode  string `json:"postcode"`
	Address   string `json:"address"`
	AddressEn string `json:"address_en"`
	Route     string `json:"route"`
	PI        bool   `json:"pi"`
	ASN       string `json:"asn"`
	V4        string `json:"v4"`
	V6        string `json:"v6"`
	V4Name    string `json:"v4_name"`
	V6Name    string `json:"v6_name"`
	Date      string `json:"date"`
	Plan      string `json:"plan"`
	Lock      bool   `json:"lock"`
}

type NetworkUser struct {
	gorm.Model
	GroupID     uint   `json:"group_id"`
	Type        uint   `json:"type"`
	Name        string `json:"name"`
	OperationID []int  `json:"operation_id"`
	TechID      []int  `json:"tech_id"`
	IP          string `json:"ip"`
	Route       string `json:"route"`
	Date        string `json:"date"`
	Plan        string `json:"plan"`
	Lock        bool   `json:"lock"`
}

type Confirm struct {
	Finish bool `json:"finish"`
}

type Result struct {
	Status  bool      `json:"status"`
	Error   string    `json:"error"`
	Network []Network `json:"network"`
}

type ResultDatabase struct {
	Err     error
	Network []Network
}
