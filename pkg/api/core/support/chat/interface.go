package chat

import "github.com/jinzhu/gorm"

const (
	ID           = 0
	TicketID     = 1
	UpdateUserID = 2
	UpdateAll    = 110
)

type Chat struct {
	gorm.Model
	TicketID uint   `json:"ticket_id"`
	UserID   uint   `json:"user_id"`
	Admin    bool   `json:"admin"`
	Data     string `json:"data"`
}

type ResultDatabase struct {
	Err  error
	Chat []Chat
}
