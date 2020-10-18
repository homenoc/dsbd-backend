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
	NetworkId uint `json:"network_id"`
	UserId    uint `json:"group_id"`
	Lock      bool `json:"lock"`
}

type Result struct {
	Status bool         `json:"status"`
	Error  string       `json:"error"`
	Jpnic  []JpnicAdmin `json:"jpnic"`
}

type ResultDatabase struct {
	Err   error
	Jpnic []JpnicAdmin
}
