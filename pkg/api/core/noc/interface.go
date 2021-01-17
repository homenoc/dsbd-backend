package noc

import (
	"github.com/jinzhu/gorm"
)

const (
	ID        = 0
	Name      = 1
	Enable    = 2
	UpdateAll = 110
)

type NOC struct {
	gorm.Model
	Name      string `json:"name"`
	Location  string `json:"location"`
	Bandwidth string `json:"bandwidth"`
	Enable    *bool  `json:"enable"`
	Comment   string `json:"comment"`
}

type Result struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	NOC    []NOC  `json:"noc"`
}

type ResultDatabase struct {
	Err error
	NOC []NOC
}
