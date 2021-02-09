package group

import (
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	jpnicAdmin "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	jpnicTech "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicTech"
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
	Agree     *bool  `json:"agree"`
	Question  string `json:"question"`
	Org       string `json:"org"`
	Status    uint   `json:"status"`
	Bandwidth string `json:"bandwidth"`
	Contract  string `json:"contract"`
	Student   *bool  `json:"student"`
	Comment   string `json:"comment"`
	Lock      *bool  `json:"lock"`
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
	Group Group `json:"group"`
}

type ResultAll struct {
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
