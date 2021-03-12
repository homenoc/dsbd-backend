package service

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID              = 0
	GID             = 1
	Org             = 2
	Open            = 3
	SearchNewNumber = 4
	UpdateName      = 100
	UpdateDate      = 102
	UpdateRoute     = 103
	UpdatePlan      = 104
	UpdateGID       = 104
	UpdateData      = 105
	UpdateAll       = 150
)

type Input struct {
	JPNICAdmin        core.JPNICAdmin  `json:"jpnic_admin"`
	JPNICTech         []core.JPNICTech `json:"jpnic_tech"`
	ServiceTemplateID uint             `json:"service_template_id"`
	ServiceComment    string           `json:"network_comment"`
	Org               string           `json:"org"`
	OrgEn             string           `json:"org_en"`
	Postcode          string           `json:"postcode"`
	Address           string           `json:"address"`
	AddressEn         string           `json:"address_en"`
	RouteV4           string           `json:"route_v4"`
	RouteV6           string           `json:"route_v6"`
	AveUpstream       uint             `json:"avg_upstream"`
	MaxUpstream       uint             `json:"max_upstream"`
	AveDownstream     uint             `json:"avg_downstream"`
	MaxDownstream     uint             `json:"max_downstream"`
	MaxBandWidthAS    uint             `json:"max_bandwidth_as"`
	PI                bool             `json:"pi"` //廃止予定
	ASN               uint             `json:"asn"`
	IP                []IPInput        `json:"ip"`
	Lock              bool             `json:"lock"`
}

type IPInput struct {
	Version   uint         `json:"version"`
	Name      string       `json:"name"`
	IP        string       `json:"ip"`
	Plan      []*core.Plan `json:"plan"`
	StartDate string       `json:"start_date"`
	EndDate   *string      `json:"end_date"`
	UseCase   string       `json:"use_case"`
}

type Confirm struct {
	Finish bool `json:"finish"`
}

type Result struct {
	Service []core.Service `json:"service"`
	//User    []core.User `json:"user"`
}

type ResultOne struct {
	Service core.Service `json:"service"`
}

type ResultDatabase struct {
	Err     error
	Service []core.Service
}
