package notice

import (
	"time"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

type Input struct {
	UserID    []uint  `json:"user_id"`
	GroupID   []uint  `json:"group_id"`
	NOCID     []uint  `json:"noc_id"`
	Everyone  *bool   `json:"everyone"`
	StartTime string  `json:"start_time"`
	EndTime   *string `json:"end_time"`
	Important *bool   `json:"important"`
	Fault     *bool   `json:"fault"`
	Info      *bool   `json:"info"`
	Title     string  `json:"title"`
	Body      string  `json:"body"`
}

type ResultAdmin struct {
	Notice []core.Notice `json:"notice"`
}

type ResultDatabase struct {
	Err    error
	Notice []core.Notice
}

// Notice is the user-facing wire shape of an active notice (GET /notice).
type Notice struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Everyone  bool      `json:"everyone"`
	Important bool      `json:"important"`
	Fault     bool      `json:"fault"`
	Info      bool      `json:"info"`
	Title     string    `json:"title"`
	Data      string    `json:"data"`
}

type Result struct {
	Notice []Notice `json:"notice"`
}
