package connection

import "github.com/homenoc/dsbd-backend/pkg/api/core"

const (
	ID        = 0
	Subnet    = 1
	UpdateAll = 150
)

type Input struct {
	Hidden  *bool  `json:"hidden"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type Result struct {
	Err  error               `json:"err"`
	IPv6 []core.IPv6Template `json:"ipv6"`
}

type ResultDatabase struct {
	Err  error
	IPv6 []core.IPv6Template
}
