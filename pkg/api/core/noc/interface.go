package noc

import (
	"github.com/jinzhu/gorm"
)

const (
	ID        = 0
	Host      = 1
	Address   = 2
	Enable    = 3
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
