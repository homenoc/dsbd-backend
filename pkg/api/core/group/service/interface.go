package service

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

type Input struct {
	JPNICAdmin     core.JPNICAdmin  `json:"jpnic_admin"`
	JPNICTech      []core.JPNICTech `json:"jpnic_tech"`
	ServiceType    string           `json:"service_type"`
	ServiceComment string           `json:"service_comment"`
	Org            string           `json:"org"`
	OrgEn          string           `json:"org_en"`
	Postcode       string           `json:"postcode"`
	Address        string           `json:"address"`
	AddressEn      string           `json:"address_en"`
	Abuse          string           `json:"abuse"`
	AveUpstream    uint             `json:"avg_upstream"`
	MaxUpstream    uint             `json:"max_upstream"`
	AveDownstream  uint             `json:"avg_downstream"`
	MaxDownstream  uint             `json:"max_downstream"`
	MaxBandWidthAS string           `json:"max_bandwidth_as"`
	StartDate      string           `json:"start_date"`
	EndDate        *string          `json:"end_date"`
	ASN            uint             `json:"asn"`
	IP             []IPInput        `json:"ip"`
	Comment        string           `json:"comment"`
	BGPComment     string           `json:"bgp_comment"`
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
