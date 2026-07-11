package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(t *core.Ticket) (*core.Ticket, error) {
	err := store.DB().Create(&t).Error
	return t, err
}

func Delete(t *core.Ticket) error {
	return store.DB().Delete(t).Error
}

// UpdateAll updates the editable ticket fields (was Update(UpdateAll, ...)).
func UpdateAll(t core.Ticket) error {
	return store.DB().Model(&core.Ticket{Model: gorm.Model{ID: t.ID}}).Updates(&core.Ticket{
		Title:         t.Title,
		GroupID:       t.GroupID,
		UserID:        t.UserID,
		Solved:        t.Solved,
		Request:       t.Request,
		RequestReject: t.RequestReject,
	}).Error
}

// GetByID loads one ticket with its user, group, and chat messages.
func GetByID(id uint) ticket.ResultDatabase {
	var tickets []core.Ticket
	err := store.DB().Preload("User").
		Preload("Group").
		Preload("Chat").
		Preload("Chat.User").
		First(&tickets, id).Error
	return ticket.ResultDatabase{Tickets: tickets, Err: err}
}

func GetAll() ticket.ResultDatabase {
	var tickets []core.Ticket
	err := store.DB().Preload("User").
		Preload("Group").
		Preload("Chat").
		Preload("Chat.User").
		Find(&tickets).Error
	return ticket.ResultDatabase{Tickets: tickets, Err: err}
}
