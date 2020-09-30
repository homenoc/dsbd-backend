package group

import "github.com/jinzhu/gorm"

const (
	ID           = 0
	OrgJa        = 1
	Org          = 2
	Email        = 3
	UpdateID     = 100
	UpdateOrg    = 101
	UpdateStatus = 102
	UpdateTechID = 103
	UpdateInfo   = 104
)

type Group struct {
	gorm.Model
	Agree     bool   `json:"agree"`
	Question  string `json:"question"`
	Org       string `json:"org"`
	Status    uint   `json:"status"`
	TechID    string `json:"tech_id"`
	Bandwidth string `json:"bandwidth"`
	Name      string `json:"name"`
	PostCode  string `json:"postcode"`
	Address   string `json:"address"`
	Mail      string `json:"mail"`
	Phone     string `json:"phone"`
	Country   string `json:"country"`
	Comment   string `json:"comment"`
}

type Result struct {
	Status    bool    `json:"status"`
	Error     string  `json:"error"`
	GroupData []Group `json:"data"`
}
