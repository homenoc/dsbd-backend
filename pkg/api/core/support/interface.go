package support

import (
	"github.com/gorilla/websocket"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/chat"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/jinzhu/gorm"
	"time"
)

// クライアントから受け取るメッセージを格納
var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan Data)

// クライアントからは JSON 形式で受け取る
type Data struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `json:"user_id"`
	Message   string    `json:"message"`
}

type FirstInput struct {
	gorm.Model
	TicketID uint   `json:"ticket_id"`
	Title    string `json:"title"`
	Data     string `json:"data"`
}

type Result struct {
	Status bool            `json:"status"`
	Error  string          `json:"error"`
	Ticket []ticket.Ticket `json:"support_ticket"`
	Chat   []chat.Chat     `json:"support_chat"`
}
