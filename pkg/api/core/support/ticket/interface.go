package ticket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

//#4 Issue(解決済み）

type Ticket struct {
	ID       uint   `json:"id"`
	Time     string `json:"time"`
	GroupID  uint   `json:"group_id"`
	UserID   uint   `json:"user_id"`
	Chat     []Chat `json:"chat"`
	Solved   *bool  `json:"solved"`
	Title    string `json:"title"`
	UserName string `json:"username"`
}

type Chat struct {
	Time     string `json:"time"`
	TicketID uint   `json:"ticket_id"`
	UserID   uint   `json:"user_id"`
	UserName string `json:"username"`
	Admin    bool   `json:"admin"`
	Data     string `json:"data"`
}

type Result struct {
	Ticket Ticket `json:"tickets"`
}

type ResultAll struct {
	Tickets []Ticket `json:"tickets"`
}

type ResultTicketAll struct {
	Tickets []Ticket `json:"tickets"`
}

type ResultAdminAll struct {
	Tickets []core.Ticket `json:"tickets"`
}

type ResultDatabase struct {
	Err     error
	Tickets []core.Ticket
}

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// UserView is the user-facing wire shape of a support ticket (GET /ticket).
type UserView struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	GroupID   uint       `json:"group_id"`
	UserID    uint       `json:"user_id"`
	Title     string     `json:"title"`
	Admin     *bool      `json:"admin"`
	Chat      []ChatView `json:"chat"`
	Solved    *bool      `json:"solved"`
}

// RequestView is the user-facing wire shape of a request ticket (GET /request).
type RequestView struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	GroupID   uint       `json:"group_id"`
	UserID    uint       `json:"user_id"`
	Title     string     `json:"title"`
	Admin     *bool      `json:"admin"`
	Chat      []ChatView `json:"chat"`
	Solved    *bool      `json:"solved"`
	Reject    *bool      `json:"reject"`
}

type ChatView struct {
	CreatedAt time.Time `json:"created_at"`
	TicketID  uint      `json:"ticket_id"`
	UserID    uint      `json:"user_id"`
	Admin     bool      `json:"admin"`
	Data      string    `json:"data" gorm:"size:65535"`
}
