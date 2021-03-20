package notice

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

const (
	ID                    = 0
	UIDOrAll              = 1
	UIDOrGIDOrAll         = 2
	UIDOrGIDOrNOCAllOrAll = 3
	NOCAll                = 4
	Important             = 10
	Fault                 = 11
	Info                  = 12
	UpdateAll             = 150
)

type Input struct {
	UserID    uint    `json:"user_id"`
	GroupID   uint    `json:"group_id"`
	NOCID     uint    `json:"noc_id"`
	Everyone  *bool   `json:"everyone"`
	StartTime string  `json:"start_time"`
	EndTime   *string `json:"end_time"`
	Important *bool   `json:"important"`
	Fault     *bool   `json:"fault"`
	Info      *bool   `json:"info"`
	Title     string  `json:"title"`
	Data      string  `json:"data"`
}

type Notice struct {
	ID        uint   `json:"ID"`
	UserID    uint   `json:"user_id"`
	GroupID   uint   `json:"group_id"`
	NOCID     uint   `json:"noc_id"`
	Everyone  bool   `json:"everyone"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Important bool   `json:"important"`
	Fault     bool   `json:"fault"`
	Info      bool   `json:"info"`
	Title     string `json:"title"`
	Data      string `json:"data" gorm:"size:65535"`
}

type Result struct {
	Notice []Notice `json:"notice"`
}

type ResultAdmin struct {
	Notice []core.Notice `json:"notice"`
}

type ResultDatabase struct {
	Err    error
	Notice []core.Notice
}
