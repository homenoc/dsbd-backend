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

//#4 Issue(解決済み）
type Ticket struct {
	gorm.Model
	GroupID     uint   `json:"group_id"`
	UserID      uint   `json:"user_id"`
	Chat        []Chat `json:"chat"`
	ChatIDStart uint   `json:"chat_id_start"`
	ChatIDEnd   uint   `json:"chat_id_end"`
	Solved      *bool  `json:"solved"`
	Title       string `json:"title"`
}

type Chat struct {
	gorm.Model
	TicketID uint   `json:"ticket_id"`
	UserID   uint   `json:"user_id"`
	Admin    bool   `json:"admin"`
	Data     string `json:"data" gorm:"size:65535"`
}

type AdminAllResult struct {
	Ticket []AdminResult `json:"ticket"`
}

type AdminResult struct {
	gorm.Model
	Status      bool   `json:"status"`
	Error       string `json:"error"`
	GroupID     uint   `json:"group_id"`
	GroupName   string `json:"group_name"`
	UserID      uint   `json:"user_id"`
	UserName    string `json:"user_name"`
	ChatIDStart uint   `json:"chat_id_start"`
	ChatIDEnd   uint   `json:"chat_id_end"`
	Solved      *bool  `json:"solved"`
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
