package notice

import "github.com/jinzhu/gorm"

type Notice struct {
	gorm.Model
	UserID     uint   `json:"user_id"`
	GroupID    uint   `json:"next_id"`
	StartTime  uint   `json:"start_time"`
	EndingTime uint   `json:"ending_time"`
	Important  bool   `json:"important"`
	Fault      bool   `json:"fault"`
	Info       bool   `json:"info"`
	Title      string `json:"title"`
	Data       string `json:"data"`
}

type Result struct {
	Status bool     `json:"status"`
	Error  string   `json:"error"`
	Notice []Notice `json:"notice"`
}

type ResultDatabase struct {
	Err    error
	Notice []Notice
}
