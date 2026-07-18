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

// Update writes the mutable ticket columns from a merged full object (both
// callers merge the request onto the current record first). Pointer columns
// are written only when provided; chats and creation metadata stay untouched.
func Update(t core.Ticket) error {
	cols := []string{"title"}
	if t.Solved != nil {
		cols = append(cols, "solved")
	}
	if t.Request != nil {
		cols = append(cols, "request")
	}
	if t.RequestReject != nil {
		cols = append(cols, "request_reject")
	}
	return store.DB().Model(&core.Ticket{Model: gorm.Model{ID: t.ID}}).Select(cols).Updates(t).Error
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
