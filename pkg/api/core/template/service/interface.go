package service

import "github.com/homenoc/dsbd-backend/pkg/api/core"

const (
	ID        = 0
	UpdateAll = 150
)

type Input struct {
	Hidden  *bool  `json:"hidden"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type Result struct {
	Err      error                  `json:"err"`
	Services []core.ServiceTemplate `json:"services"`
}

type ResultDatabase struct {
	Err      error
	Services []core.ServiceTemplate
}
