package noc

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID        = 0
	Name      = 1
	Enable    = 2
	UpdateAll = 150
)

type ResultOneUser struct {
	ID       uint   `json:"ID"`
	Name     string `json:"name"`
	Location string `json:"location"`
	New      *bool  `json:"new"`
}

type ResultAllUser struct {
	NOC []ResultOneUser `json:"noc"`
}

type Result struct {
	NOC []core.NOC `json:"noc"`
}

type ResultDatabase struct {
	Err error
	NOC []core.NOC
}
