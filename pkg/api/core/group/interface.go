package group

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID           = 0
	OrgJa        = 1
	Org          = 2
	Email        = 3
	UpdateID     = 100
	UpdateOrg    = 101
	UpdateStatus = 102
	UpdateTechID = 103
	UpdateInfo   = 104
	UpdateAll    = 150
)

type Input struct {
	Agree          *bool   `json:"agree"`
	Question       string  `json:"question"`
	Org            string  `json:"org"`
	Contract       string  `json:"contract"`
	Student        *bool   `json:"student"`
	StudentExpired *string `json:"student_expired"`
}

type Result struct {
	Status bool         `json:"status"`
	Error  string       `json:"error"`
	Group  []core.Group `json:"group"`
}

type AdminResult struct {
	//User       []core.User             `json:"user"`
	Group      []core.Group      `json:"group"`
	Network    []core.Service    `json:"network"`
	Connection []core.Connection `json:"connection"`
}

type ResultOne struct {
	ID            uint   `json:"id"`
	Agree         *bool  `json:"agree"`
	Question      string `json:"question"`
	Org           string `json:"org"`
	OrgEn         string `json:"org_en"`
	PostCode      string `json:"postcode"`
	Address       string `json:"address"`
	AddressEn     string `json:"address_en"`
	Tel           string `json:"tel"`
	Country       string `json:"country"`
	Status        uint   `json:"status"`
	Bandwidth     string `json:"bandwidth"`
	Contract      string `json:"contract"`
	Student       *bool  `json:"student"`
	Pass          *bool  `json:"pass"`
	Lock          *bool  `json:"lock"`
	ExpiredStatus uint   `json:"expired_status"`
	Open          *bool  `json:"open"`
}

type ResultAll struct {
	Group   ResultOne      `json:"group"`
	Network []core.Service `json:"network"`
	//Admin      []core.Admin      `json:"admin"`
	//Tech       []core.Tech       `json:"tech"`
	Connection []core.Connection `json:"connection"`
}

type ResultDatabase struct {
	Err   error
	Group []core.Group
}
