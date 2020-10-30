package user

import "github.com/jinzhu/gorm"

const (
	ID               = 0
	GID              = 1
	Name             = 2
	Email            = 3
	MailToken        = 4
	UpdateVerifyMail = 100
	UpdateGID        = 101
	UpdateInfo       = 102
	UpdateStatus     = 105
	UpdateLevel      = 106
	UpdateAll        = 110
)

type User struct {
	gorm.Model
	GID        uint   `json:"gid"`
	Tech       *bool  `json:"tech"`
	Name       string `json:"name"`
	NameEn     string `json:"name_en"`
	Email      string `json:"email"`
	Pass       string `json:"pass"`
	Status     uint   `json:"status"`
	Level      uint   `json:"level"`
	MailVerify *bool  `json:"mail_verify"`
	MailToken  string `json:"mail_token"`
	Org        string `json:"org"`
	OrgEn      string `json:"org_en"`
	PostCode   string `json:"postcode"`
	Address    string `json:"address"`
	AddressEn  string `json:"address_en"`
	Dept       string `json:"dept"`
	DeptEn     string `json:"dept_en"`
	Pos        string `json:"pos"`
	PosEn      string `json:"pos_en"`
	Tel        string `json:"tel"`
	Fax        string `json:"fax"`
	Country    string `json:"country"`
}

type ResultOne struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	User   User   `json:"data"`
}

type Result struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	User   []User `json:"data"`
}

type ResultDatabase struct {
	Err  error
	User []User
}
