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

// UpdateAll updates all mutable user fields (was Update(UpdateAll, ...)).
func UpdateAll(u *core.User) error {
	return store.DB().Model(&core.User{Model: gorm.Model{ID: u.ID}}).Updates(core.User{
		GroupID:           u.GroupID,
		Name:              u.Name,
		NameEn:            u.NameEn,
		Email:             u.Email,
		Pass:              u.Pass,
		Level:             u.Level,
		MailVerify:        u.MailVerify,
		MailToken:         u.MailToken,
		ExpiredStatus:     u.ExpiredStatus,
		AntisocialCheck:   u.AntisocialCheck,
		AntisocialCheckAt: u.AntisocialCheckAt,
	}).Error
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
