package chat

import "github.com/jinzhu/gorm"

const (
	ID        = 0
	TicketID  = 1
	UpdateAll = 110
)

type Chat struct {
	gorm.Model
	TicketID uint   `json:"ticket_id"`
	NextID   uint   `json:"next_id"`
	Admin    bool   `json:"admin"`
	Data     string `json:"data"`
}

type ResultDatabase struct {
	Err  error
	Chat []Chat
}
