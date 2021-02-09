package jpnicAdmin

import "github.com/jinzhu/gorm"

const (
	ID               = 0
	NetworkId        = 1
	UserId           = 2
	NetworkAndUserId = 3
	UpdateLock       = 100
	UpdateAll        = 110
)

type JpnicAdmin struct {
	gorm.Model
	NetworkID uint  `json:"network_id"`
	UserID    uint  `json:"user_id"`
	Lock      *bool `json:"lock"`
}

type Result struct {
	Jpnic []JpnicAdmin `json:"jpnic"`
}

type ResultDatabase struct {
	Err   error
	Jpnic []JpnicAdmin
}
