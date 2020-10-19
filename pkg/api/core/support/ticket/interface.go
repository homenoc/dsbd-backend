package ticket

import "github.com/jinzhu/gorm"

const (
	ID        = 0
	GID       = 1
	UID       = 2
	CID       = 2
	UpdateAll = 110
)

type Ticket struct {
	gorm.Model
	GroupID uint   `json:"group_id"`
	UserID  uint   `json:"user_id"`
	ChatID  uint   `json:"chat_id"`
	Title   string `json:"title"`
}

type ResultDatabase struct {
	Err    error
	Ticket []Ticket
}
