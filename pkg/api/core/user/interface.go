package user

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/info"
)

const (
	ID               = 0
	IDDetail         = 10
	GID              = 1
	Name             = 2
	Email            = 3
	MailToken        = 4
	GIDAndLevel      = 5
	IDGetGroup       = 6
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

type User struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	NameEn        string `json:"name_en"`
	Email         string `json:"email"`
	Level         uint   `json:"level"`
	ExpiredStatus uint   `json:"expired_status"`
	MailVerify    *bool  `json:"mail_verify"`
}

type SimpleGroup struct {
	ID            uint  `json:"id"`
	Student       *bool `json:"student"`
	Pass          *bool `json:"pass"`
	Lock          *bool `json:"lock"`
	ExpiredStatus *uint `json:"expired_status"`
	Status        *uint `json:"status"`
}

type ResultOne struct {
	User  SimpleUser  `json:"user"`
	Group SimpleGroup `json:"group"`
	Info  []info.Info `json:"info"`
}

type SimpleUser struct {
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
