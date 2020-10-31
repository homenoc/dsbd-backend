package group

import (
	"github.com/jinzhu/gorm"
)

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
	PI        *bool  `json:"pi"`
	ASN       string `json:"asn"`
	RouteV4   string `json:"route_v4"`
	RouteV6   string `json:"route_v6"`
	V4        string `json:"v4"`
	V6        string `json:"v6"`
	V4Name    string `json:"v4_name"`
	V6Name    string `json:"v6_name"`
	Date      string `json:"date"`
	Plan      string `json:"plan"`
	Open      *bool  `json:"open"`
	Lock      *bool  `json:"lock"`
}

type NetworkInput struct {
	AdminID   uint   `json:"admin_id"`
	TechID    []uint `json:"tech_id"`
	GroupID   uint   `json:"group_id"`
	Org       string `json:"org"`
	OrgEn     string `json:"org_en"`
	Postcode  string `json:"postcode"`
	Address   string `json:"address"`
	AddressEn string `json:"address_en"`
	RouteV4   string `json:"route_v4"`
	RouteV6   string `json:"route_v6"`
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

type Confirm struct {
	Finish bool `json:"finish"`
}

type Result struct {
	Status  bool      `json:"status"`
	Error   string    `json:"error"`
	Network []Network `json:"network"`
}

type ResultOne struct {
	Status  bool    `json:"status"`
	Error   string  `json:"error"`
	Network Network `json:"network"`
}

type ResultDatabase struct {
	Err     error
	Network []Network
}
