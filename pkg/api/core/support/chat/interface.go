package chat

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID           = 0
	TicketID     = 1
	UpdateUserID = 2
	UpdateAll    = 150
)

type ResultDatabase struct {
	Err  error
	Chat []core.Chat
}
