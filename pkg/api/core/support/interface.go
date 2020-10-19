package support

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/chat"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/jinzhu/gorm"
)

type FirstInput struct {
	gorm.Model
	Title string `json:"title"`
	Data  string `json:"data"`
}

type Result struct {
	Status bool            `json:"status"`
	Error  string          `json:"error"`
	Ticket []ticket.Ticket `json:"support_ticket"`
	Chat   []chat.Chat     `json:"support_chat"`
}
