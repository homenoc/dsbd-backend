package admin

import (
	"github.com/jinzhu/gorm"
)

const (
	ID         = 0
	UserID     = 10
	Name       = 11
	UpdateName = 100
	UpdateAll  = 101
)

type Admin struct {
	gorm.Model
	UserID   uint   `json:"user_id"`
	FullName string `json:"full_name"`
	NickName string `json:"nick_name"`
	Active   bool   `json:"active"`
	Debug    string `json:"debug"`
}

type Result struct {
	Status bool    `json:"status"`
	Error  string  `json:"error"`
	Admin  []Admin `json:"admin"`
}

type ResultTmpToken struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Admin  string `json:"admin"`
}

type ResultDatabase struct {
	Err   error
	Admin []Admin
}
