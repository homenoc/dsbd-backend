package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/chat"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
)

func Create(s *core.Chat) (*core.Chat, error) {
	err := store.DB().Create(&s).Error
	return s, err
}

func Delete(s *core.Chat) error {
	return store.DB().Delete(s).Error
}

// GetByID looks up one chat message by primary key.
func GetByID(id uint) chat.ResultDatabase {
	var chats []core.Chat
	err := store.DB().First(&chats, id).Error
	return chat.ResultDatabase{Chat: chats, Err: err}
}

// GetByTicketID returns a ticket's chat messages in ascending order.
func GetByTicketID(ticketID uint) chat.ResultDatabase {
	var chats []core.Chat
	err := store.DB().Where("ticket_id = ?", ticketID).Order("id asc").Find(&chats).Error
	return chat.ResultDatabase{Chat: chats, Err: err}
}

func GetAll() chat.ResultDatabase {
	var chats []core.Chat
	err := store.DB().Find(&chats).Error
	return chat.ResultDatabase{Chat: chats, Err: err}
}
