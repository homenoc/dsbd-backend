package core

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Tokens        []*Token  `json:"tokens"`
	Notice        []*Notice `json:"notice"`
	GroupID       uint      `json:"group_id"`
	Name          string    `json:"name"`
	NameEn        string    `json:"name_en"`
	Email         string    `json:"email"`
	Pass          string    `json:"pass"`
	ExpiredStatus *uint     `json:"expired_status"`
	Level         uint      `json:"level"`
	MailVerify    *bool     `json:"mail_verify"`
	MailToken     string    `json:"mail_token"`
}

type Group struct {
	gorm.Model
	Users          []User        `json:"users"`
	Services       []Service     `json:"services"`
	Tickets        []Ticket      `json:"tickets"`
	Notice         []*Notice     `json:"notice"`
	JPNICAdmin     []*JPNICAdmin `json:"jpnic_admin"`
	JPNICTech      []*JPNICTech  `json:"jpnic_tech"`
	Agree          *bool         `json:"agree"`
	Question       string        `json:"question"  gorm:"size:65535"`
	Org            string        `json:"org"`
	OrgEn          string        `json:"org_en"`
	PostCode       string        `json:"postcode"`
	Address        string        `json:"address"`
	AddressEn      string        `json:"address_en"`
	Tel            string        `json:"tel"`
	Country        string        `json:"country"`
	Status         *uint         `json:"status"`
	Contract       string        `json:"contract"`
	Student        *bool         `json:"student"`
	StudentExpired *time.Time    `json:"student_expired"`
	Fee            *uint         `json:"fee"`
	Comment        string        `json:"comment"`
	Pass           *bool         `json:"pass"`
	Lock           *bool         `json:"lock"`
	ExpiredStatus  *uint         `json:"expired_status"`
}

type Service struct {
	gorm.Model
	GroupID           uint             `json:"group_id"`
	ServiceTemplateID *uint            `json:"service_template_id"`
	ServiceTemplate   *ServiceTemplate `json:"service_template"`
	ServiceComment    string           `json:"service_comment"`
	ServiceNumber     uint             `json:"service_number"`
	Org               string           `json:"org"`
	OrgEn             string           `json:"org_en"`
	Postcode          string           `json:"postcode"`
	Address           string           `json:"address"`
	AddressEn         string           `json:"address_en"`
	ASN               uint             `json:"asn"`
	RouteV4           string           `json:"route_v4"`
	RouteV6           string           `json:"route_v6"`
	V4Name            string           `json:"v4_name"`
	V6Name            string           `json:"v6_name"`
	AveUpstream       uint             `json:"avg_upstream"`
	MaxUpstream       uint             `json:"max_upstream"`
	AveDownstream     uint             `json:"avg_downstream"`
	MaxDownstream     uint             `json:"max_downstream"`
	MaxBandWidthAS    uint             `json:"max_bandwidth_as"`
	Fee               *uint            `json:"fee"`
	IP                []IP             `json:"ip"`
	Connections       []Connection     `json:"connections"`
	JPNICAdminID      uint             `json:"jpnic_admin_id"`
	JPNICAdmin        JPNICAdmin       `json:"jpnic_admin"`
	JPNICTech         []JPNICTech      `json:"jpnic_tech" gorm:"many2many:service_jpnic_tech;"`
	Open              *bool            `json:"open"`
	Lock              *bool            `json:"lock"`
	AddAllow          *bool            `json:"add_allow"`
}

type Connection struct {
	gorm.Model
	ServiceID                uint                   `json:"service_id"`
	BGPRouterID              *uint                  `json:"bgp_router_id"`                //使用RouterのID
	TunnelEndPointRouterIPID *uint                  `json:"tunnel_endpoint_router_ip_id"` //使用エンドポイントルータのID
	ConnectionTemplateID     *uint                  `json:"connection_template_id"`
	ConnectionTemplate       *ConnectionTemplate    `json:"connection_template"`
	ConnectionComment        string                 `json:"connection_comment"` // ServiceがETCの時や補足説明で必要
	ConnectionNumber         uint                   `json:"connection_number"`
	NTTTemplateID            *uint                  `json:"ntt_template_id"`
	NOCID                    *uint                  `json:"noc_id"`
	TermIP                   string                 `json:"term_ip"`
	Monitor                  *bool                  `json:"monitor"`
	Address                  string                 `json:"address"` //都道府県　市町村
	LinkV4Our                string                 `json:"link_v4_our"`
	LinkV4Your               string                 `json:"link_v4_your"`
	LinkV6Our                string                 `json:"link_v6_our"`
	LinkV6Your               string                 `json:"link_v6_your"`
	Open                     *bool                  `json:"open"`
	Lock                     *bool                  `json:"lock"`
	Comment                  string                 `json:"comment"`
	NTTTemplate              *NTTTemplate           `json:"ntt_template"`
	NOC                      *NOC                   `json:"noc"`
	BGPRouter                BGPRouter              `json:"bgp_router"`
	TunnelEndPointRouterIP   TunnelEndPointRouterIP `json:"tunnel_endpoint_router_ip"`
}

type NOC struct {
	gorm.Model
	Notice               []*Notice               `json:"notice"`
	BGPRouter            []*BGPRouter            `json:"bgp_router"`
	TunnelEndPointRouter []*TunnelEndPointRouter `json:"tunnel_endpoint_router"`
	Name                 string                  `json:"name"`
	Location             string                  `json:"location"`
	Bandwidth            string                  `json:"bandwidth"`
	Enable               *bool                   `json:"enable"`
	Comment              string                  `json:"comment"`
}

type BGPRouter struct {
	gorm.Model
	NOCID        uint                    `json:"noc_id"`
	HostName     string                  `json:"hostname"`
	Address      string                  `json:"address"`
	TunnelRouter []*TunnelEndPointRouter `json:"tunnel_endpoint_router"`
	Enable       *bool                   `json:"enable"`
	Comment      string                  `json:"comment"`
}

type TunnelEndPointRouter struct {
	gorm.Model
	NOCID                  uint                      `json:"noc_id"`
	TunnelEndPointRouterIP []*TunnelEndPointRouterIP `json:"tunnel_endpoint_router_ip"`
	HostName               string                    `json:"hostname"`
	Capacity               uint                      `json:"capacity"`
	Enable                 *bool                     `json:"enable"`
	Comment                string                    `json:"comment"`
}

type TunnelEndPointRouterIP struct {
	gorm.Model
	TunnelRouterID uint   `json:"tunnel_router_id"`
	IP             string `json:"ip"`
	Enable         *bool  `json:"enable"`
	Comment        string `json:"comment"`
}

type IP struct {
	gorm.Model
	ServiceID uint       `json:"service_id"`
	Version   uint       `json:"version"`
	Name      string     `json:"name"`
	IP        string     `json:"ip"`
	Plan      []*Plan    `json:"plan" `
	PlanJPNIC *string    `json:"" gorm:"size:65535"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	UseCase   string     `json:"use_case"`
	Open      *bool      `json:"open"`
}

type Plan struct {
	gorm.Model
	IPID     uint   `json:"ip_id"`
	Name     string `json:"name"`
	After    uint   `json:"after"`
	HalfYear uint   `json:"half_year"`
	OneYear  uint   `json:"one_year"`
}

type JPNICAdmin struct {
	gorm.Model
	Service     []Service `gorm:"foreignkey:JPNICAdminID"`
	GroupName   string    `json:"group_name"`
	GroupNameEn string    `json:"group_name_en"`
	Org         string    `json:"org"`
	OrgEn       string    `json:"org_en"`
	PostCode    string    `json:"postcode"`
	Address     string    `json:"address"`
	AddressEn   string    `json:"address_en"`
	Dept        string    `json:"dept"`
	DeptEn      string    `json:"dept_en"`
	Pos         string    `json:"pos"`
	PosEn       string    `json:"pos_en"`
	Tel         string    `json:"tel"`
	Fax         string    `json:"fax"`
	Country     string    `json:"country"`
	Lock        *bool     `json:"lock"`
}

type JPNICTech struct {
	gorm.Model
	Service     []Service `json:"service" gorm:"many2many:service_jpnic_tech;"`
	GroupName   string    `json:"group_name"`
	GroupNameEn string    `json:"group_name_en"`
	Org         string    `json:"org"`
	OrgEn       string    `json:"org_en"`
	PostCode    string    `json:"postcode"`
	Address     string    `json:"address"`
	AddressEn   string    `json:"address_en"`
	Dept        string    `json:"dept"`
	DeptEn      string    `json:"dept_en"`
	Pos         string    `json:"pos"`
	PosEn       string    `json:"pos_en"`
	Tel         string    `json:"tel"`
	Fax         string    `json:"fax"`
	Country     string    `json:"country"`
	Lock        *bool     `json:"lock"`
}

type ServiceTemplate struct {
	gorm.Model
	Hidden       *bool  `json:"hidden"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Comment      string `json:"comment"`
	NeedJPNIC    *bool  `json:"need_jpnic"`
	NeedGlobalAS *bool  `json:"need_global_as"`
	NeedComment  *bool  `json:"need_comment"`
	NeedRoute    *bool  `json:"need_route"`
}

type ConnectionTemplate struct {
	gorm.Model
	Hidden           bool   `json:"hidden"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Comment          string `json:"comment"`
	NeedInternet     *bool  `json:"need_internet"`
	NeedComment      *bool  `json:"need_comment"`
	NeedCrossConnect *bool  `json:"need_cross_connect"`
}

type NTTTemplate struct {
	gorm.Model
	Hidden  bool   `json:"hidden"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type Ticket struct {
	gorm.Model
	GroupID uint   `json:"group_id"`
	UserID  uint   `json:"user_id"`
	Chat    []Chat `json:"chat"`
	Solved  *bool  `json:"solved"`
	Title   string `json:"title"`
}

type Chat struct {
	gorm.Model
	TicketID uint   `json:"ticket_id"`
	UserID   uint   `json:"user_id"`
	Admin    bool   `json:"admin"`
	Data     string `json:"data" gorm:"size:65535"`
}

type Token struct {
	gorm.Model
	ExpiredAt   time.Time `json:"expired_at"`
	UserID      uint      `json:"user_id"`
	Status      uint      `json:"status"` //0: initToken(30m) 1: 30m 2:6h 3: 12h 10: 30d 11:180d
	Admin       *bool     `json:"admin"`
	UserToken   string    `json:"user_token"`
	TmpToken    string    `json:"tmp_token"`
	AccessToken string    `json:"access_token"`
	Debug       string    `json:"debug"`
}

type Notice struct {
	gorm.Model
	UserID     uint      `json:"user_id"`
	GroupID    uint      `json:"group_id"`
	Everyone   *bool     `json:"everyone"`
	StartTime  time.Time `json:"start_time"`
	EndingTime time.Time `json:"ending_time"`
	Important  *bool     `json:"important"`
	Fault      *bool     `json:"fault"`
	Info       *bool     `json:"info"`
	Title      string    `json:"title"`
	Data       string    `json:"data"`
}
