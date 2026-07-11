package user

import (
	"time"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

type Input struct {
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
	Email  string `json:"email"`
	Pass   string `json:"pass"`
	Level  uint   `json:"level"`
}

// Profile is the user-facing wire shape of an account (GET /user/me and the
// group member list in GET /group).
type Profile struct {
	ID                uint       `json:"id"`
	GroupID           uint       `json:"group_id"`
	StripeCustomerID  string     `json:"stripe_customer_id"`
	Name              string     `json:"name"`
	NameEn            string     `json:"name_en"`
	Email             string     `json:"email"`
	Status            uint       `json:"status"`
	Level             uint       `json:"level"`
	MailVerify        *bool      `json:"mail_verify"`
	AntisocialCheck   *bool      `json:"antisocial_check"`
	AntisocialCheckAt *time.Time `json:"antisocial_check_at"`
}

// ProfileFrom projects a core.User onto its user-facing wire shape.
func ProfileFrom(u core.User) Profile {
	var groupID uint
	if u.GroupID != nil {
		groupID = *u.GroupID
	}
	return Profile{
		ID:                u.ID,
		GroupID:           groupID,
		Name:              u.Name,
		NameEn:            u.NameEn,
		Email:             u.Email,
		Level:             u.Level,
		MailVerify:        u.MailVerify,
		AntisocialCheck:   u.AntisocialCheck,
		AntisocialCheckAt: u.AntisocialCheckAt,
	}
}

type Result struct {
	User []Profile `json:"user"`
}

type ResultAdmin struct {
	User []core.User `json:"users"`
}

type ResultDatabase struct {
	Err  error
	User []core.User
}
