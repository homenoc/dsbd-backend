package group

import (
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	jpnicAdmin "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	jpnicTech "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicTech"
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
	Agree     bool   `json:"agree"`
	Question  string `json:"question"`
	Org       string `json:"org"`
	Status    uint   `json:"status"`
	Bandwidth string `json:"bandwidth"`
	Contract  string `json:"contract"`
	Comment   string `json:"comment"`
	Lock      bool   `json:"lock"`
}

type Result struct {
	Status bool    `json:"status"`
	Error  string  `json:"error"`
	Group  []Group `json:"group"`
}

type ResultOne struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Group  Group  `json:"data"`
}

type ResultAll struct {
	Status     bool                    `json:"status"`
	Error      string                  `json:"error"`
	Group      Group                   `json:"group"`
	Network    []network.Network       `json:"network"`
	JpnicAdmin []jpnicAdmin.JpnicAdmin `json:"admin"`
	JpnicTech  []jpnicTech.JpnicTech   `json:"tech"`
	Connection []connection.Connection `json:"connection"`
}

type ResultDatabase struct {
	Err   error
	Group []Group
}
