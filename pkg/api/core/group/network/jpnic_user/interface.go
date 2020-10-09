package jpnic_user

import "github.com/jinzhu/gorm"

const (
	ID           = 0
	Name         = 1
	Mail         = 2
	GID          = 3
	UpdateID     = 100
	UpdateOpID   = 101
	UpdateTechID = 102
	UpdateGID    = 103
	UpdateInfo   = 104
	UpdateAll    = 110
)

type JPNICUser struct {
	gorm.Model
	GroupID     uint   `json:"group_id"`
	OperationID uint   `json:"operation_id"`
	TechID      uint   `json:"tech_id"`
	NameJa      string `json:"name_ja"`
	Name        string `json:"name"`
	OrgJa       string `json:"org_ja"`
	Org         string `json:"org"`
	PostCode    string `json:"postcode"`
	AddressJa   string `json:"address_ja"`
	Address     string `json:"address"`
	DeptJa      string `json:"dept_ja"`
	Dept        string `json:"dept"`
	PosJa       string `json:"pos_ja"`
	Pos         string `json:"pos"`
	Mail        string `json:"mail"`
	Tel         string `json:"tel"`
	Fax         string `json:"fax"`
}

type Result struct {
	Status    bool        `json:"status"`
	Error     string      `json:"error"`
	JPNICUser []JPNICUser `json:"jpnic_user"`
}

type ResultDatabase struct {
	Err       error
	JPNICUser []JPNICUser
}
