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
	RouterID uint   `json:"router_id"`
	HostName string `json:"hostname"`
	V4       string `json:"v4"`
	V6       string `json:"v6"`
	Capacity uint   `json:"capacity"`
	Enable   bool   `json:"enable"`
}

type Result struct {
	Gateway []Gateway `json:"gateway"`
}

type ResultDatabase struct {
	Err     error
	Gateway []Gateway
}
