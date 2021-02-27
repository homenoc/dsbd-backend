package network

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	ID          = 0
	GID         = 1
	Org         = 2
	Open        = 3
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
	GroupID     uint                    `json:"group_id"`
	NetworkType string                  `json:"network_type"`
	Org         string                  `json:"org"`
	OrgEn       string                  `json:"org_en"`
	Postcode    string                  `json:"postcode"`
	Address     string                  `json:"address"`
	AddressEn   string                  `json:"address_en"`
	PI          *bool                   `json:"pi"`
	ASN         string                  `json:"asn"`
	RouteV4     string                  `json:"route_v4"`
	RouteV6     string                  `json:"route_v6"`
	V4Name      string                  `json:"v4_name"`
	V6Name      string                  `json:"v6_name"`
	IP          []IP                    `json:"ip"`
	Connection  []connection.Connection `json:"connection"`
	JPNICAdmin  JPNICAdmin              `json:"jpnic_admin"`
	JPNICTech   []JPNICTech             `json:"jpnic_tech"`
	Plan        string                  `json:"plan"`
	Open        *bool                   `json:"open"`
	Lock        *bool                   `json:"lock"`
}

type IP struct {
	gorm.Model
	NetworkID uint       `json:"network_id"`
	Version   uint       `json:"version"`
	Name      string     `json:"name"`
	IP        string     `json:"ip"`
	Plan      *string    `json:"plan"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	UseCase   string     `json:"use_case"`
	Open      *bool      `json:"open"`
}

type JPNICAdmin struct {
	gorm.Model
	NetworkID uint  `json:"network_id"`
	UserID    uint  `json:"user_id"`
	Lock      *bool `json:"lock"`
}

type JPNICTech struct {
	gorm.Model
	NetworkID uint  `json:"network_id"`
	UserID    uint  `json:"user_id"`
	Lock      *bool `json:"lock"`
}

type Input struct {
	AdminID   uint       `json:"admin_id"`
	TechID    []uint     `json:"tech_id"`
	GroupID   uint       `json:"group_id"`
	Org       string     `json:"org"`
	OrgEn     string     `json:"org_en"`
	Postcode  string     `json:"postcode"`
	Address   string     `json:"address"`
	AddressEn string     `json:"address_en"`
	RouteV4   string     `json:"route_v4"`
	RouteV6   string     `json:"route_v6"`
	PI        bool       `json:"pi"`
	ASN       string     `json:"asn"`
	IP        *[]IPInput `json:"ip"`
	Lock      bool       `json:"lock"`
}

type IPInput struct {
	Version   uint    `json:"version"`
	Name      string  `json:"name"`
	IP        string  `json:"ip"`
	Plan      *string `json:"plan"`
	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date"`
	UseCase   string  `json:"use_case"`
}

type Confirm struct {
	Finish bool `json:"finish"`
}

type Result struct {
	Network []Network   `json:"network"`
	User    []user.User `json:"user"`
}

type ResultOne struct {
	Network Network `json:"network"`
}

type ResultDatabase struct {
	Err     error
	Network []Network
}
