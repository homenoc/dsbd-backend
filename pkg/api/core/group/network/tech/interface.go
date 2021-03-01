package tech

import "github.com/jinzhu/gorm"

const (
	ID               = 0
	NetworkID        = 1
	UserId           = 2
	NetworkAndUserId = 3
	UpdateLock       = 100
	UpdateAll        = 110
)

type Tech struct {
	gorm.Model
	NetworkID uint  `json:"network_id"`
	UserID    uint  `json:"user_id"`
	Lock      *bool `json:"lock"`
}

type Result struct {
	Tech []Tech `json:"tech"`
}

type ResultDatabase struct {
	Err  error
	Tech []Tech
}
