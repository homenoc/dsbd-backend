package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(u *core.User) error {
	return store.DB().Create(&u).Error
}

func Delete(u *core.User) error {
	return store.DB().Delete(u).Error
}

// UpdateMailVerify sets only the mail-verify flag (was Update(UpdateVerifyMail, ...)).
func UpdateMailVerify(u *core.User) error {
	return store.DB().Model(&core.User{Model: gorm.Model{ID: u.ID}}).
		Updates(core.User{MailVerify: u.MailVerify}).Error
}

// UpdateGID sets only the group id (was Update(UpdateGID, ...)).
func UpdateGID(u *core.User) error {
	return store.DB().Model(&core.User{Model: gorm.Model{ID: u.ID}}).
		Updates(core.User{GroupID: u.GroupID}).Error
}

// UpdateAntisocialCheck sets the antisocial-check flag and timestamp.
func UpdateAntisocialCheck(u *core.User) error {
	return store.DB().Model(&core.User{Model: gorm.Model{ID: u.ID}}).
		Updates(core.User{AntisocialCheck: u.AntisocialCheck, AntisocialCheckAt: u.AntisocialCheckAt}).Error
}

// Update writes the editable user columns from a merged full object (both
// callers merge the request onto the current record first). Credential-ish
// columns (pass, mail_token) and level are written only when non-empty/non-zero
// — they are never legitimately cleared; pointer columns only when provided.
func Update(u *core.User) error {
	cols := []string{"name", "name_en", "email"}
	if u.Pass != "" {
		cols = append(cols, "pass")
	}
	if u.MailToken != "" {
		cols = append(cols, "mail_token")
	}
	if u.Level != 0 {
		cols = append(cols, "level")
	}
	if u.GroupID != nil {
		cols = append(cols, "group_id")
	}
	if u.MailVerify != nil {
		cols = append(cols, "mail_verify")
	}
	if u.ExpiredStatus != nil {
		cols = append(cols, "expired_status")
	}
	if u.AntisocialCheck != nil {
		cols = append(cols, "antisocial_check")
	}
	if u.AntisocialCheckAt != nil {
		cols = append(cols, "antisocial_check_at")
	}
	return store.DB().Model(&core.User{Model: gorm.Model{ID: u.ID}}).Select(cols).Updates(u).Error
}

// GetByID looks up one user by primary key.
func GetByID(id uint) user.ResultDatabase {
	var users []core.User
	err := store.DB().First(&users, id).Error
	return user.ResultDatabase{User: users, Err: err}
}

// GetDetail loads a user with the full group/service/connection association graph.
func GetDetail(id uint) user.ResultDatabase {
	var users []core.User
	err := store.DB().Where("id = ?", id).
		Preload("Ticket").
		Preload("Ticket.Chat").
		Preload("Group").
		Preload("Group.Users").
		Preload("Group.Services").
		Preload("Group.Tickets").
		Preload("Group.Tickets.Chat").
		Preload("Group.Services.IP").
		Preload("Group.Services.IP.Plan").
		Preload("Group.Services.Connection").
		Preload("Group.Services.Connection.BGPRouter").
		Preload("Group.Services.Connection.BGPRouter.NOC").
		Preload("Group.Services.Connection.TunnelEndPointRouterIP").
		Preload("Group.Services.JPNICAdmin").
		Preload("Group.Services.JPNICTech").Find(&users).Error
	return user.ResultDatabase{User: users, Err: err}
}

// GetByEmail returns users matching an email address.
func GetByEmail(email string) user.ResultDatabase {
	var users []core.User
	err := store.DB().Where("email = ?", email).Find(&users).Error
	return user.ResultDatabase{User: users, Err: err}
}

// GetByMailToken returns users matching an email-verification token.
func GetByMailToken(mailToken string) user.ResultDatabase {
	var users []core.User
	err := store.DB().Where("mail_token = ?", mailToken).Find(&users).Error
	return user.ResultDatabase{User: users, Err: err}
}

// GetByGroupIDAndLevel returns users in a group at a given level.
func GetByGroupIDAndLevel(groupID *uint, level uint) user.ResultDatabase {
	var users []core.User
	err := store.DB().Where("group_id = ? AND level = ?", groupID, level).Find(&users).Error
	return user.ResultDatabase{User: users, Err: err}
}

func GetArray(u []uint) user.ResultDatabase {
	var users []core.User
	err := store.DB().Where(u).Find(&users).Error
	return user.ResultDatabase{User: users, Err: err}
}

func GetAll() user.ResultDatabase {
	var users []core.User
	err := store.DB().Find(&users).Error
	return user.ResultDatabase{User: users, Err: err}
}
