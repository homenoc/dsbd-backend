package token

import (
	"github.com/jinzhu/gorm"
	"time"
)

const (
	ID                      = 0
	UserToken               = 10
	UserTokenAndAccessToken = 11
	ExpiredTime             = 12
	AdminToken              = 20
	AddToken                = 100
	UpdateToken             = 101
	UpdateAll               = 110
)

type Token struct {
	gorm.Model
	ExpiredAt   time.Time `json:"expired_at"`
	UserID      uint      `json:"user_id"`
	AdminID     uint      `json:"admin_id"`
	Status      uint      `json:"status"` //0: initToken(30m) 1: 30m 2:6h 3: 12h 10: 30d 11:180d
	Admin       *bool     `json:"admin"`
	UserToken   string    `json:"user_token"`
	TmpToken    string    `json:"tmp_token"`
	AccessToken string    `json:"access_token"`
	Debug       string    `json:"debug"`
}

type Result struct {
	Status bool    `json:"status"`
	Error  string  `json:"error"`
	Token  []Token `json:"token"`
}

type ResultTmpToken struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Token  string `json:"token"`
}

type ResultDatabase struct {
	Err   error
	Token []Token
}
