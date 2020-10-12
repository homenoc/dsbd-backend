package network_user

import "github.com/jinzhu/gorm"

const (
	ID         = 0
	NetAndType = 1
	Network    = 2
	User       = 3
	Type       = 4
	Info       = 5
	UpdateInfo = 100
	UpdateAll  = 110
)

type NetworkUser struct {
	gorm.Model
	Type        uint `json:"type"` //1: op 2: tech
	NetworkID   uint `json:"network_id"`
	JPNICUserID uint `json:"jpnic_user_id"`
	Lock        bool `json:"lock"`
}

type Result struct {
	Status      bool          `json:"status"`
	Error       string        `json:"error"`
	NetworkUser []NetworkUser `json:"network"`
}

type ResultDatabase struct {
	Err         error
	NetworkUser []NetworkUser
}
