package info

import (
	"time"
)

type User struct {
	ID               uint   `json:"id"`
	GroupID          uint   `json:"group_id"`
	StripeCustomerID string `json:"stripe_customer_id"`
	Name             string `json:"name"`
	NameEn           string `json:"name_en"`
	Email            string `json:"email"`
	Status           uint   `json:"status"`
	Level            uint   `json:"level"`
	MailVerify       *bool  `json:"mail_verify"`
}

type Group struct {
	ID            uint       `json:"id"`
	Agree         *bool      `json:"agree"`
	Question      string     `json:"question"`
	Org           string     `json:"org"`
	OrgEn         string     `json:"org_en"`
	PostCode      string     `json:"postcode"`
	Address       string     `json:"address"`
	AddressEn     string     `json:"address_en"`
	Tel           string     `json:"tel"`
	Country       string     `json:"country"`
	Contract      string     `json:"contract"`
	CouponID      string     `json:"coupon_id"`
	MemberTypeID  uint       `json:"member_type_id"`
	MemberType    string     `json:"member_type"`
	MemberExpired *time.Time `json:"member_expired"`
	IsExpired     bool       `json:"is_expired"`
	IsStripeID    bool       `json:"is_stripe_id"`
	Pass          *bool      `json:"pass"`
	ExpiredStatus *uint      `json:"expired_status"`
	AddAllow      *bool      `json:"add_allow"`
}

type Notice struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Everyone  bool      `json:"everyone"`
	Important bool      `json:"important"`
	Fault     bool      `json:"fault"`
	Info      bool      `json:"info"`
	Title     string    `json:"title"`
	Data      string    `json:"data"`
}

type Ticket struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	GroupID   uint      `json:"group_id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Admin     *bool     `json:"admin"`
	Chat      []Chat    `json:"chat"`
	Solved    *bool     `json:"solved"`
}

type Chat struct {
	CreatedAt time.Time `json:"created_at"`
	TicketID  uint      `json:"ticket_id"`
	UserID    uint      `json:"user_id"`
	Admin     bool      `json:"admin"`
	Data      string    `json:"data" gorm:"size:65535"`
}

type Info struct {
	ServiceID      string   `json:"service_id"`
	Service        string   `json:"service"`
	Assign         bool     `json:"assign"`
	ASN            uint     `json:"asn"`
	V4             []string `json:"v4"`
	V6             []string `json:"v6"`
	NOC            string   `json:"noc"`
	NOCIP          string   `json:"noc_ip"`
	TermIP         string   `json:"term_ip"`
	LinkV4Our      string   `json:"link_v4_our"`
	LinkV4Your     string   `json:"link_v4_your"`
	LinkV6Our      string   `json:"link_v6_our"`
	LinkV6Your     string   `json:"link_v6_your"`
	Fee            string   `json:"fee"`
	Org            string   `json:"org"`
	OrgEn          string   `json:"org_en"`
	PostCode       string   `json:"postcode"`
	Address        string   `json:"address"`
	AddressEn      string   `json:"address_en"`
	JPNICAdmin     JPNIC    `json:"jpnic_admin"`
	JPNICTech      []JPNIC  `json:"jpnic_tech"`
	AveUpstream    uint     `json:"avg_upstream"`
	MaxUpstream    uint     `json:"max_upstream"`
	AveDownstream  uint     `json:"avg_downstream"`
	MaxDownstream  uint     `json:"max_downstream"`
	MaxBandWidthAS string   `json:"max_bandwidth_as"`
	BGPRouteV4     string   `json:"bgp_route_v4"`
	BGPRouteV6     string   `json:"bgp_route_v6"`
}

type JPNIC struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	NameEn    string `json:"name_en"`
	Mail      string `json:"mail"`
	Org       string `json:"org"`
	OrgEn     string `json:"org_en"`
	PostCode  string `json:"postcode"`
	Address   string `json:"address"`
	AddressEn string `json:"address_en"`
	Dept      string `json:"dept"`
	DeptEn    string `json:"dept_en"`
	Tel       string `json:"tel"`
	Fax       string `json:"fax"`
	Country   string `json:"country"`
}

type IP struct {
	ID        uint   `json:"id"`
	Version   uint   `json:"version"`
	Name      string `json:"name"`
	IP        string `json:"ip"`
	Plan      []Plan `json:"plan" `
	PlanJPNIC string `json:"" gorm:"size:65535"`
	UseCase   string `json:"use_case"`
}

type Plan struct {
	ID       uint   `json:"id"`
	IPID     uint   `json:"ip_id"`
	Name     string `json:"name"`
	After    uint   `json:"after"`
	HalfYear uint   `json:"half_year"`
	OneYear  uint   `json:"one_year"`
}

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

type Connection struct {
	ID        uint   `json:"id"`
	ServiceID string `json:"service_id"`
	Open      bool   `json:"open"`
}

type Request struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	GroupID   uint      `json:"group_id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Admin     *bool     `json:"admin"`
	Chat      []Chat    `json:"chat"`
	Solved    *bool     `json:"solved"`
	Reject    *bool     `json:"reject"`
}

type Result struct {
	User       User         `json:"user"`
	Group      Group        `json:"group"`
	UserList   []User       `json:"user_list"`
	Notice     []Notice     `json:"notice"`
	Ticket     []Ticket     `json:"ticket"`
	Service    []Service    `json:"service"`
	Connection []Connection `json:"connection"`
	Request    []Request    `json:"request"`
	Info       []Info       `json:"info"`
}
