package gateway

import (
	"github.com/jinzhu/gorm"
)

const (
	ID        = 0
	NOC       = 1
	Address   = 2
	Enable    = 3
	UpdateAll = 110
)

type Gateway struct {
	gorm.Model
	NOCID    uint   `json:"noc_id"`
	HostName string `json:"hostname"`
	Capacity uint   `json:"capacity"`
	Comment  string `json:"comment"`
	Enable   bool   `json:"enable"`
}

type Result struct {
	Gateway []Gateway `json:"gateway"`
}

type ResultDatabase struct {
	Err     error
	Gateway []Gateway
}
