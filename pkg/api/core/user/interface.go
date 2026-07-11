package user

import (
	"time"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/info"
)

type Input struct {
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
	Email  string `json:"email"`
	Pass   string `json:"pass"`
	Level  uint   `json:"level"`
}

type User struct {
	ID                uint       `json:"id"`
	Name              string     `json:"name"`
	NameEn            string     `json:"name_en"`
	Email             string     `json:"email"`
	Level             uint       `json:"level"`
	ExpiredStatus     uint       `json:"expired_status"`
	MailVerify        *bool      `json:"mail_verify"`
	AntisocialCheck   *bool      `json:"antisocial_check"`
	AntisocialCheckAt *time.Time `json:"antisocial_check_at"`
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
	ID                uint       `json:"id"`
	GroupID           uint       `json:"group_id"`
	Name              string     `json:"name"`
	NameEn            string     `json:"name_en"`
	Email             string     `json:"email"`
	Status            uint       `json:"status"`
	Level             uint       `json:"level"`
	MailVerify        *bool      `json:"mail_verify"`
	AntisocialCheck   *bool      `json:"antisocial_check"`
	AntisocialCheckAt *time.Time `json:"antisocial_check_at"`
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
