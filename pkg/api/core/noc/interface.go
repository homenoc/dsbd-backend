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
	New       *bool  `json:"new"`
	Enable    *bool  `json:"enable"`
	Comment   string `json:"comment"`
}

type ResultOneUser struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	New      *bool  `json:"new"`
}

type ResultAllUser struct {
	NOC []ResultOneUser `json:"noc"`
}

type Result struct {
	NOC []NOC `json:"noc"`
}

type ResultDatabase struct {
	Err error
	NOC []NOC
}
