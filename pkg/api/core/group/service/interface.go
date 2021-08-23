package service

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
)

const (
	ID                = 0
	GID               = 1
	Org               = 2
	Open              = 3
	SearchNewNumber   = 4
	GIDAndAddAllow    = 5
	UpdateName        = 100
	UpdateDate        = 102
	UpdateRoute       = 103
	UpdateStatus      = 104
	UpdateGID         = 104
	UpdateData        = 105
	ReplaceJPNICAdmin = 120
	ReplaceJPNICTech  = 121
	ReplaceConnection = 122
	ReplaceIP         = 123
	UpdateAll         = 199
	AppendJPNICAdmin  = 200
	AppendJPNICTech   = 201
	AppendConnection  = 202
	AppendIP          = 203
	DeleteJPNICAdmin  = 300
	DeleteJPNICTech   = 301
	DeleteConnection  = 302
	DeleteIP          = 303
)

type Input struct {
	JPNICAdmin        core.JPNICAdmin  `json:"jpnic_admin"`
	JPNICTech         []core.JPNICTech `json:"jpnic_tech"`
	ServiceTemplateID uint             `json:"service_template_id"`
	ServiceComment    string           `json:"service_comment"`
	Org               string           `json:"org"`
	OrgEn             string           `json:"org_en"`
	Postcode          string           `json:"postcode"`
	Address           string           `json:"address"`
	AddressEn         string           `json:"address_en"`
	AveUpstream       uint             `json:"avg_upstream"`
	MaxUpstream       uint             `json:"max_upstream"`
	AveDownstream     uint             `json:"avg_downstream"`
	MaxDownstream     uint             `json:"max_downstream"`
	MaxBandWidthAS    string           `json:"max_bandwidth_as"`
	StartDate         string           `json:"start_date"`
	EndDate           *string          `json:"end_date"`
	ASN               uint             `json:"asn"`
	IP                []IPInput        `json:"ip"`
	Lock              bool             `json:"lock"`
}

type JPNIC struct {
	ID          uint   `json:"id"`
	JPNICHandle string `json:"jpnic_handle"`
	Name        string `json:"name"`
	NameEn      string `json:"name_en"`
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

type Service struct {
	ID                  uint                     `json:"id"`
	GroupID             uint                     `json:"group_id"`
	ServiceTemplateID   *uint                    `json:"service_template_id"`
	ServiceTemplateName string                   `json:"service_template_name"`
	ServiceComment      string                   `json:"service_comment"`
	ServiceNumber       uint                     `json:"service_number"`
	Org                 string                   `json:"org"`
	OrgEn               string                   `json:"org_en"`
	PostCode            string                   `json:"postcode"`
	Address             string                   `json:"address"`
	AddressEn           string                   `json:"address_en"`
	ASN                 *uint                    `json:"asn"`
	RouteV4             string                   `json:"route_v4"`
	RouteV6             string                   `json:"route_v6"`
	V4Name              string                   `json:"v4_name"`
	V6Name              string                   `json:"v6_name"`
	AveUpstream         uint                     `json:"avg_upstream"`
	MaxUpstream         uint                     `json:"max_upstream"`
	AveDownstream       uint                     `json:"avg_downstream"`
	MaxDownstream       uint                     `json:"max_downstream"`
	MaxBandWidthAS      string                   `json:"max_bandwidth_as"`
	Fee                 *uint                    `json:"fee"`
	IP                  []core.IP                `json:"ip"`
	Connections         *[]connection.Connection `json:"connections"`
	JPNICAdminID        uint                     `json:"jpnic_admin_id"`
	JPNICAdmin          *JPNIC                   `json:"jpnic_admin"`
	JPNICTech           *[]JPNIC                 `json:"jpnic_tech"`
	Open                *bool                    `json:"open"`
	AddAllow            *bool                    `json:"add_allow"`
	Lock                *bool                    `json:"lock"`
}

type IP struct {
	ID        uint         `json:"id"`
	Version   uint         `json:"version"`
	Name      string       `json:"name"`
	IP        string       `json:"ip"`
	Plan      []*core.Plan `json:"plan"`
	StartDate string       `json:"start_date"`
	EndDate   *string      `json:"end_date"`
	UseCase   string       `json:"use_case"`
	Open      *bool        `json:"open"`
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
