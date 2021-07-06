package core

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Tokens        []*Token  `json:"tokens"`
	Notice        []*Notice `json:"notice" gorm:"many2many:user_notice;"`
	Ticket        []Ticket  `json:"tickets"`
	Group         *Group    `json:"group"`
	Payment       []Payment `json:"payment_membership"`
	GroupID       *uint     `json:"group_id"`
	Name          string    `json:"name"`
	NameEn        string    `json:"name_en"`
	Email         string    `json:"email"`
	Pass          string    `json:"pass"`
	ExpiredStatus *uint     `json:"expired_status"`
	Level         uint      `json:"level"`
	MailVerify    *bool     `json:"mail_verify"`
	MailToken     string    `json:"mail_token"`
}

type Payment struct {
	gorm.Model
	User            *User  `json:"user"`
	Group           *Group `json:"group"`
	UserID          uint   `json:"user_id"`
	GroupID         *uint  `json:"group_id"`
	PaymentIntentID string `json:"payment_intent_id"`
	IsMembership    *bool  `json:"is_membership"`
	Paid            *bool  `json:"paid"`
	Refund          *bool  `json:"refund"`
	Fee             uint   `json:"fee"`
	Comment         string `json:"comment"`
}

type Group struct {
	gorm.Model
	Users                       []User                    `json:"users"`
	Payment                     []Payment                 `json:"payment_membership"`
	Services                    []Service                 `json:"services"`
	Tickets                     []Ticket                  `json:"tickets"`
	PaymentMembershipTemplateID *uint                     `json:"payment_membership_template_id"`
	PaymentCouponTemplateID     *uint                     `json:"payment_coupon_template_id"`
	PaymentMembershipTemplate   PaymentMembershipTemplate `json:"payment_membership_template"`
	PaymentCouponTemplate       PaymentCouponTemplate     `json:"payment_coupon_template"`
	StripeCustomerID            *string                   `json:"stripe_customer_id"`
	StripePaymentMethodID       *string                   `json:"stripe_payment_method_id"` //Todo: いらんかも
	StripeSubscriptionID        *string                   `json:"stripe_subscription_id"`
	Agree                       *bool                     `json:"agree"`
	Question                    string                    `json:"question"  gorm:"size:10000"`
	Org                         string                    `json:"org"`
	OrgEn                       string                    `json:"org_en"`
	PostCode                    string                    `json:"postcode"`
	Address                     string                    `json:"address"`
	AddressEn                   string                    `json:"address_en"`
	Tel                         string                    `json:"tel"`
	Country                     string                    `json:"country"`
	Status                      *uint                     `json:"status"`
	Contract                    string                    `json:"contract"`
	MemberExpired               *time.Time                `json:"member_expired"`
	Student                     *bool                     `json:"student"`
	Fee                         *uint                     `json:"fee"`
	Comment                     string                    `json:"comment"`
	Open                        *bool                     `json:"open"`
	Pass                        *bool                     `json:"pass"`
	Lock                        *bool                     `json:"lock"` //いらん
	ExpiredStatus               *uint                     `json:"expired_status"`
	AddAllow                    *bool                     `json:"add_allow"`
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
	PostCode          string           `json:"postcode"`
	Address           string           `json:"address"`
	AddressEn         string           `json:"address_en"`
	ASN               *uint            `json:"asn"`
	RouteV4           string           `json:"route_v4"`
	RouteV6           string           `json:"route_v6"`
	V4Name            string           `json:"v4_name"`
	V6Name            string           `json:"v6_name"`
	AveUpstream       uint             `json:"avg_upstream"`
	MaxUpstream       uint             `json:"max_upstream"`
	AveDownstream     uint             `json:"avg_downstream"`
	MaxDownstream     uint             `json:"max_downstream"`
	MaxBandWidthAS    string           `json:"max_bandwidth_as"`
	Fee               *uint            `json:"fee"` //いらんかも
	IP                []IP             `json:"ip"`
	Connection        []*Connection    `json:"connections"`
	JPNICAdmin        JPNICAdmin       `json:"jpnic_admin"`
	JPNICTech         []JPNICTech      `json:"jpnic_tech"`
	StartDate         time.Time        `json:"start_date"`
	EndDate           *time.Time       `json:"end_date"`
	Pass              *bool            `json:"pass"`
	Enable            *bool            `json:"enable"`
	Lock              *bool            `json:"lock"` //いらん
	AddAllow          *bool            `json:"add_allow"`
	Group             Group            `json:"group"`
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
	IPv4RouteTemplateID      *uint                  `json:"ipv4_route_template_id"`
	IPv6RouteTemplateID      *uint                  `json:"ipv6_route_template_id"`
	NOCID                    *uint                  `json:"noc_id"`
	TermIP                   string                 `json:"term_ip"`
	Monitor                  *bool                  `json:"monitor"`
	Address                  string                 `json:"address"` //都道府県　市町村
	LinkV4Our                string                 `json:"link_v4_our"`
	LinkV4Your               string                 `json:"link_v4_your"`
	LinkV6Our                string                 `json:"link_v6_our"`
	LinkV6Your               string                 `json:"link_v6_your"`
	Open                     *bool                  `json:"open"`
	Enable                   *bool                  `json:"enable"`
	Lock                     *bool                  `json:"lock"` //いらん
	Comment                  string                 `json:"comment"`
	IPv4RouteTemplate        *IPv4RouteTemplate     `json:"ipv4_route_template"`
	IPv6RouteTemplate        *IPv6RouteTemplate     `json:"ipv6_route_template"`
	NTTTemplate              *NTTTemplate           `json:"ntt_template"`
	NOC                      *NOC                   `json:"noc"`
	BGPRouter                BGPRouter              `json:"bgp_router"`
	TunnelEndPointRouterIP   TunnelEndPointRouterIP `json:"tunnel_endpoint_router_ip"`
	Service                  Service                `json:"service"`
}

type NOC struct {
	gorm.Model
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
	NOCID uint `json:"noc_id"`
	NOC   NOC  `json:"noc"`
	//Connection []Connection `json:"connection"`
	HostName string `json:"hostname"`
	Address  string `json:"address"`
	Enable   *bool  `json:"enable"`
	Comment  string `json:"comment"`
}

type TunnelEndPointRouter struct {
	gorm.Model
	NOCID                  *uint                     `json:"noc_id"`
	TunnelEndPointRouterIP []*TunnelEndPointRouterIP `json:"tunnel_endpoint_router_ip"`
	HostName               string                    `json:"hostname"`
	Capacity               uint                      `json:"capacity"`
	Enable                 *bool                     `json:"enable"`
	Comment                string                    `json:"comment"`
}

type TunnelEndPointRouterIP struct {
	gorm.Model
	TunnelEndPointRouter   TunnelEndPointRouter `json:"tunnel_endpoint_router"`
	TunnelEndPointRouterID *uint                `json:"tunnel_endpoint_router_id"`
	IP                     string               `json:"ip"`
	Enable                 *bool                `json:"enable"`
	Comment                string               `json:"comment"`
}

type IP struct {
	gorm.Model
	ServiceID uint       `json:"service_id"`
	Version   uint       `json:"version"`
	Name      string     `json:"name"`
	IP        string     `json:"ip"`
	Plan      []*Plan    `json:"plan" `
	PlanJPNIC *string    `json:"" gorm:"size:15000"` //いらんかも
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
	ServiceID uint `json:"service_id"`
	//Service     []Service `gorm:"foreignkey:JPNICAdminID"`
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
	Lock        *bool  `json:"lock"`
}

type JPNICTech struct {
	gorm.Model
	//Service     []Service `json:"service" gorm:"many2many:service_jpnic_tech;"`
	ServiceID   uint   `json:"service_id"`
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
	Lock        *bool  `json:"lock"`
}

type PaymentMembershipTemplate struct {
	gorm.Model
	PriceID string `json:"price_id"`
	Title   string `json:"title"`
	Plan    string `json:"plan"`
	Monthly bool   `json:"monthly"`
	Yearly  bool   `json:"yearly"`
	Fee     uint   `json:"fee"`
	Comment string `json:"comment"`
}

type PaymentCouponTemplate struct {
	gorm.Model
	StripeCouponID string `json:"stripe_coupon_id"`
	Title          string `json:"title"`
	DiscountRate   uint   `json:"discount_rate"`
	Comment        string `json:"comment"`
}

type PaymentDonateTemplate struct {
	gorm.Model
	Name    string `json:"name"`
	Fee     uint   `json:"fee"`
	Comment string `json:"comment"`
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

type IPv4RouteTemplate struct {
	gorm.Model
	Name string `json:"name"`
}

type IPv6RouteTemplate struct {
	gorm.Model
	Name string `json:"name"`
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
	L2               *bool  `json:"l2"`
	L3               *bool  `json:"l3"`
}

type NTTTemplate struct {
	gorm.Model
	Hidden  bool   `json:"hidden"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// 申請中/承諾済み/却下
type Ticket struct {
	gorm.Model
	GroupID       *uint  `json:"group_id"`
	UserID        *uint  `json:"user_id"`
	Chat          []Chat `json:"chat"`
	Request       *bool  `json:"request"`
	RequestReject *bool  `json:"request_reject"`
	Solved        *bool  `json:"solved"`
	Admin         *bool  `json:"admin"`
	Title         string `json:"title"`
	Group         Group  `json:"group"`
	User          User   `json:"user"`
}

type Chat struct {
	gorm.Model
	TicketID uint   `json:"ticket_id"`
	UserID   *uint  `json:"user_id"`
	Admin    bool   `json:"admin"`
	Data     string `json:"data" gorm:"size:10000"`
	User     User   `json:"user"`
}

type Token struct {
	gorm.Model
	ExpiredAt   time.Time `json:"expired_at"`
	UserID      *uint     `json:"user_id"`
	User        User      `json:"user"`
	Status      uint      `json:"status"` //0: initToken(30m) 1: 30m 2:6h 3: 12h 10: 30d 11:180d
	Admin       *bool     `json:"admin"`
	UserToken   string    `json:"user_token"`
	TmpToken    string    `json:"tmp_token"`
	AccessToken string    `json:"access_token"`
	Debug       string    `json:"debug"`
}

type Notice struct {
	gorm.Model
	User      []User    `json:"user" gorm:"many2many:notice_user;"`
	Everyone  *bool     `json:"everyone"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Important *bool     `json:"important"`
	Fault     *bool     `json:"fault"`
	Info      *bool     `json:"info"`
	Title     string    `json:"title"`
	Data      string    `json:"data" gorm:"size:15000"`
}

type Request struct {
	gorm.Model
	RequestTemplateID uint            `json:"request_template_id"`
	RequestTemplate   RequestTemplate `json:"request_template"`
	TargetID          *uint           `json:"target_id"`
	Reason1           string          `json:"reason_1"`
	Reason2           string          `json:"reason_2"`
	Accept            *bool           `json:"accept"`
	User              User            `json:"user"`
	Group             Group           `json:"group"`
}

type IPv4Template struct {
	gorm.Model
	Title    string `json:"title"`
	Subnet   string `json:"subnet"`
	Quantity uint   `json:"quantity"`
	Hide     *bool  `json:"hide"`
	Comment  string `json:"comment"`
}

type IPv6Template struct {
	gorm.Model
	Title    string `json:"title"`
	Subnet   string `json:"subnet"`
	Quantity uint   `json:"quantity"`
	Hide     *bool  `json:"hide"`
	Comment  string `json:"comment"`
}

// Type 1:追加 2:修正 3:削除
// InfoType 1:グループ情報 2:サービス情報 3:IP 4:JPNICAdmin 5:JPNICTech 6:接続情報
type RequestTemplate struct {
	gorm.Model
	Title       string `json:"title"`
	Data        string `json:"data"`
	RequestType uint   `json:"request_type"`
	InfoType    uint   `json:"info_type"`
	Comment     string `json:"comment"`
}
