package group

import "github.com/jinzhu/gorm"

const (
	ID          = 0
	GID         = 1
	Name        = 2
	Type        = 3
	UpdateName  = 100
	UpdateIP    = 101
	UpdateDate  = 102
	UpdateRoute = 103
	UpdatePlan  = 104
	UpdateGID   = 104
	UpdateData  = 105
	UpdateAll   = 110
)

type Network struct {
	gorm.Model `json:"base"`
	GroupID    uint   `json:"group_id"`
	Type       uint   `json:"type"`
	Name       string `json:"name"`
	IP         string `json:"ip"`
	Route      string `json:"route"`
	Date       string `json:"date"`
	Plan       string `json:"plan"`
	Lock       bool   `json:"lock"`
}

type NetworkUser struct {
	gorm.Model  `json:"base"`
	GroupID     uint   `json:"group_id"`
	Type        uint   `json:"type"`
	Name        string `json:"name"`
	OperationID []int  `json:"operation_id"`
	TechID      []int  `json:"tech_id"`
	IP          string `json:"ip"`
	Route       string `json:"route"`
	Date        string `json:"date"`
	Plan        string `json:"plan"`
	Lock        bool   `json:"lock"`
}

type NetworkJPNICUser struct {
	gorm.Model  `json:"base"`
	NetworkID   uint `json:"network_id"`
	JPNICUserID uint `json:"jpnic_user_id"`
}

type Confirm struct {
	Finish bool `json:"finish"`
}

type Result struct {
	Status  bool      `json:"status"`
	Error   string    `json:"error"`
	Network []Network `json:"network"`
}

type ResultDatabase struct {
	Err     error
	Network []Network
}
