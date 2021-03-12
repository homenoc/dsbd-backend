package ticket

import (
	"github.com/gorilla/websocket"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/jinzhu/gorm"
	"net/http"
)

const (
	ID        = 0
	GID       = 1
	UID       = 2
	UpdateAll = 150
)

//#4 Issue(解決済み）

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
	Ticket []core.Ticket
}

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
