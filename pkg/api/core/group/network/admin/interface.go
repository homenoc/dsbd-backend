package admin

import "github.com/jinzhu/gorm"

const (
	ID               = 0
	NetworkId        = 1
	UserId           = 2
	NetworkAndUserId = 3
	UpdateLock       = 100
	UpdateAll        = 110
)

type Admin struct {
	gorm.Model
	NetworkID uint  `json:"network_id"`
	UserID    uint  `json:"user_id"`
	Lock      *bool `json:"lock"`
}

type Result struct {
	Admins []Admin `json:"admins"`
}

type ResultDatabase struct {
	Err    error
	Admins []Admin
}
