package noc

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

type Router struct {
	gorm.Model
	NOC      uint   `json:"noc"`
	HostName string `json:"hostname"`
	Address  string `json:"address"`
	Enable   *bool  `json:"enable"`
}

type Result struct {
	Status bool     `json:"status"`
	Error  string   `json:"error"`
	Router []Router `json:"router"`
}

type ResultDatabase struct {
	Err    error
	Router []Router
}
