package ticket

import (
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"net/http"
)

const (
	ID        = 0
	GID       = 1
	UID       = 2
	CIDStart  = 3
	CIDEnd    = 4
	UpdateAll = 110
)

type Ticket struct {
	gorm.Model
	GroupID     uint   `json:"group_id"`
	UserID      uint   `json:"user_id"`
	ChatIDStart uint   `json:"chat_id_start"`
	ChatIDEnd   uint   `json:"chat_id_end"`
	Solved      bool   `json:"solved"`
	Title       string `json:"title"`
}

type ResultDatabase struct {
	Err    error
	Ticket []Ticket
}

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
