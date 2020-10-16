package group

import (
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	jpnicUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnic_user"
	networkUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/network_user"
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
	Status    bool    `json:"status"`
	Error     string  `json:"error"`
	GroupData []Group `json:"data"`
}

type ResultAll struct {
	Status      bool                      `json:"status"`
	Error       string                    `json:"error"`
	Group       Group                     `json:"group"`
	Network     []network.Network         `json:"network"`
	JPNICUser   []jpnicUser.JPNICUser     `json:"jpnic_user"`
	NetworkUser []networkUser.NetworkUser `json:"network_user"`
	Connection  []connection.Connection   `json:"connection"`
}

type ResultDatabase struct {
	Err   error
	Group []Group
}
