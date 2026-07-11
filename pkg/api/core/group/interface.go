package group

import (
	"time"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

type Input struct {
	Agree          *bool   `json:"agree"`
	Question       string  `json:"question"`
	Org            string  `json:"org"`
	OrgEn          string  `json:"org_en"`
	PostCode       string  `json:"postcode"`
	Address        string  `json:"address"`
	AddressEn      string  `json:"address_en"`
	Tel            string  `json:"tel"`
	Country        string  `json:"country"`
	Contract       string  `json:"contract"`
	Student        *bool   `json:"student"`
	StudentExpired *string `json:"student_expired"`
}

type ResultAdmin struct {
	Group core.Group `json:"group"`
}

type ResultAdminAll struct {
	Group []core.Group `json:"group"`
}

type ResultDatabase struct {
	Err   error
	Group []core.Group
}

// Profile is the user-facing wire shape of the caller's own group
// (GET /group), including membership/billing status. Address-level fields are
// filled only for levels that may manage the group.
type Profile struct {
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
