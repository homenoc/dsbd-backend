package group

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/admin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/tech"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/jinzhu/gorm"
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
	UpdateAll    = 110
)

type Group struct {
	gorm.Model
	Agree         *bool  `json:"agree"`
	Question      string `json:"question"`
	Org           string `json:"org"`
	Status        *uint  `json:"status"`
	Bandwidth     string `json:"bandwidth"`
	Contract      string `json:"contract"`
	Student       *bool  `json:"student"`
	Fee           *uint  `json:"fee"`
	Comment       string `json:"comment"`
	Pass          *bool  `json:"pass"`
	Lock          *bool  `json:"lock"`
	ExpiredStatus *uint  `json:"expired_status"`
}

type Result struct {
	Status bool    `json:"status"`
	Error  string  `json:"error"`
	Group  []Group `json:"group"`
}

type AdminResult struct {
	User       []user.User             `json:"user"`
	Group      []Group                 `json:"group"`
	Network    []network.Network       `json:"network"`
	Connection []connection.Connection `json:"connection"`
}

type ResultOne struct {
	ID            uint   `json:"id"`
	Agree         *bool  `json:"agree"`
	Question      string `json:"question"`
	Org           string `json:"org"`
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
	Group      ResultOne               `json:"group"`
	Network    []network.Network       `json:"network"`
	Admin      []admin.Admin           `json:"admin"`
	Tech       []tech.Tech             `json:"tech"`
	Connection []connection.Connection `json:"connection"`
}

type ResultDatabase struct {
	Err   error
	Group []Group
}
