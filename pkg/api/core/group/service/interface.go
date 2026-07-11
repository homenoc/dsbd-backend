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

// JPNIC is a JPNIC contact as shown to users. The admin contact omits
// its address fields on the wire.
type JPNIC struct {
	ID          uint   `json:"id"`
	JPNICHandle string `json:"jpnic_handle"`
	Name        string `json:"name"`
	NameEn      string `json:"name_en"`
	Mail        string `json:"mail"`
	Org         string `json:"org"`
	OrgEn       string `json:"org_en"`
	PostCode    string `json:"postcode"`
	Address     string `json:"address"`
	AddressEn   string `json:"address_en"`
	Dept        string `json:"dept"`
	DeptEn      string `json:"dept_en"`
	Tel         string `json:"tel"`
	Fax         string `json:"fax"`
	Country     string `json:"country"`
}

// Service is the user-facing wire shape of an enabled service
// (GET /service), including the capability flags from the type registry.
type Service struct {
	ID             uint    `json:"id"`
	ServiceID      string  `json:"service_id"`
	ServiceType    string  `json:"service_type"`
	NeedRoute      bool    `json:"need_route"`
	NeedBGP        bool    `json:"need_bgp"`
	NeedJPNIC      bool    `json:"need_jpnic"`
	AddAllow       bool    `json:"add_allow"`
	Pass           bool    `json:"pass"`
	Org            string  `json:"org"`
	OrgEn          string  `json:"org_en"`
	PostCode       string  `json:"postcode"`
	Address        string  `json:"address"`
	AddressEn      string  `json:"address_en"`
	ASN            *uint   `json:"asn"`
	AveUpstream    uint    `json:"avg_upstream"`
	MaxUpstream    uint    `json:"max_upstream"`
	AveDownstream  uint    `json:"avg_downstream"`
	MaxDownstream  uint    `json:"max_downstream"`
	MaxBandWidthAS string  `json:"max_bandwidth_as"`
	JPNICAdmin     JPNIC   `json:"jpnic_admin"`
	JPNICTech      []JPNIC `json:"jpnic_tech"`
	IP             []IP    `json:"ip"`
}

// IP is an opened IP assignment as shown to users.
type IP struct {
	ID        uint   `json:"id"`
	Version   uint   `json:"version"`
	Name      string `json:"name"`
	IP        string `json:"ip"`
	Plan      []Plan `json:"plan" `
	PlanJPNIC string `json:"" gorm:"size:65535"`
	UseCase   string `json:"use_case"`
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

type Plan struct {
	ID       uint   `json:"id"`
	IPID     uint   `json:"ip_id"`
	Name     string `json:"name"`
	After    uint   `json:"after"`
	HalfYear uint   `json:"half_year"`
	OneYear  uint   `json:"one_year"`
}
