package user

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID               = 0
	GID              = 1
	Name             = 2
	Email            = 3
	MailToken        = 4
	GIDAndLevel      = 5
	UpdateVerifyMail = 100
	UpdateGID        = 101
	UpdateInfo       = 102
	UpdateStatus     = 105
	UpdateLevel      = 106
	UpdateAll        = 150
)

type Input struct {
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
	Email  string `json:"email"`
	Pass   string `json:"pass"`
	Level  uint   `json:"level"`
}

type ResultOne struct {
	ID         uint   `json:"id"`
	GroupID    uint   `json:"group_id"`
	Name       string `json:"name"`
	NameEn     string `json:"name_en"`
	Email      string `json:"email"`
	Status     uint   `json:"status"`
	Level      uint   `json:"level"`
	MailVerify *bool  `json:"mail_verify"`
}

type Result struct {
	User []ResultOne `json:"user"`
}

type ResultAdmin struct {
	User []core.User `json:"users"`
}

type ResultDatabase struct {
	Err  error
	User []core.User
}
