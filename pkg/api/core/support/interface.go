package support

import (
	"github.com/gorilla/websocket"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/chat"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/jinzhu/gorm"
	"time"
)

// channel定義(websocketで使用)
var Clients = make(map[*WebSocket]bool)
var Broadcast = make(chan WebSocketResult)

// websocket用
type WebSocketResult struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UserToken   string    `json:"user_token"`
	AccessToken string    `json:"access_token"`
	UserID      uint      `json:"user_id"`
	Message     string    `json:"message"`
}

type WebSocket struct {
	TicketID uint
	GroupID  uint
	UserID   uint
	Socket   *websocket.Conn
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
