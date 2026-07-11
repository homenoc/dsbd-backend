package group

import (
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
