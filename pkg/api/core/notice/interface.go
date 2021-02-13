package notice

import (
	"github.com/jinzhu/gorm"
)

const (
	ID               = 0
	UserID           = 1
	GroupID          = 2
	UserIDAndGroupID = 3
	Everyone         = 4
	GroupData        = 5
	UserData         = 6
	Important        = 10
	Fault            = 11
	Info             = 12
	UpdateAll        = 110
)

type Notice struct {
	gorm.Model
	UserID     uint   `json:"user_id"`
	GroupID    uint   `json:"group_id"`
	Everyone   *bool  `json:"everyone"`
	StartTime  uint   `json:"start_time"`
	EndingTime uint   `json:"ending_time"`
	Important  *bool  `json:"important"`
	Fault      *bool  `json:"fault"`
	Info       *bool  `json:"info"`
	Title      string `json:"title"`
	Data       string `json:"data"`
}

type Result struct {
	Notice []Notice `json:"notice"`
}

type ResultDatabase struct {
	Err    error
	Notice []Notice
}
