package support

import (
	"github.com/gorilla/websocket"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/jinzhu/gorm"
	"time"
)

// channel定義(websocketで使用)
var Clients = make(map[*WebSocket]bool)
var Broadcast = make(chan WebSocketResult)

// websocket用
type WebSocketResult struct {
	ID          uint      `json:"id"`
	Err         string    `json:"error"`
	CreatedAt   time.Time `json:"created_at"`
	UserToken   string    `json:"user_token"`
	AccessToken string    `json:"access_token"`
	UserID      uint      `json:"user_id"`
	UserName    string    `json:"user_name"`
	GroupID     uint      `json:"group_id"`
	Admin       bool      `json:"admin"`
	Message     string    `json:"message"`
}

type WebSocketChatResponse struct {
	Time     string `json:"time"`
	UserID   uint   `json:"user_id"`
	UserName string `json:"username"`
	GroupID  uint   `json:"group_id"`
	Admin    bool   `json:"admin"`
	Message  string `json:"message"`
}

type WebSocket struct {
	TicketID uint
	GroupID  uint
	UserID   uint
	UserName string
	Admin    bool
	Socket   *websocket.Conn
}

type FirstInput struct {
	gorm.Model
	IsGroup  bool   `json:"is_group"`
	TicketID uint   `json:"ticket_id"`
	Title    string `json:"title"`
	Data     string `json:"data"`
	UserID   uint   `json:"user_id"`
	GroupID  uint   `json:"group_id"`
}

type Result struct {
	Ticket []core.Ticket `json:"ticket"`
	Chat   []core.Chat   `json:"chat"`
}
