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
	UpdateMail       = 103
	UpdatePass       = 104
	UpdateStatus     = 105
	UpdateLevel      = 106
	UpdateAll        = 110
)

type User struct {
	gorm.Model
	GID        uint   `json:"gid"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Pass       string `json:"pass"`
	Status     uint   `json:"status"`
	Level      uint   `json:"level"`
	MailVerify bool   `json:"mail_verify"`
	MailToken  string `json:"mail_token"`
}

type Result struct {
	Status   bool   `json:"status"`
	Error    string `json:"error"`
	UserData []User `json:"data"`
}

type ResultDatabase struct {
	Err  error
	User []User
}
