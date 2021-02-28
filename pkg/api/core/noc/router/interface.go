package noc

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/gateway"
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
	NOC      uint              `json:"noc"`
	HostName string            `json:"hostname"`
	Address  string            `json:"address"`
	Gateway  []gateway.Gateway `json:"gateway"`
	Enable   *bool             `json:"enable"`
}

type Result struct {
	Router []Router `json:"router"`
}

type ResultDatabase struct {
	Err    error
	Router []Router
}
