package ticket

import (
	"github.com/gorilla/websocket"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"net/http"
)

const (
	ID        = 0
	GID       = 1
	UID       = 2
	UpdateAll = 150
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
