package noc

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
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
