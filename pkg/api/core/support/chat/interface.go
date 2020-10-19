package chat

import "github.com/jinzhu/gorm"

const (
	ID        = 0
	UpdateAll = 110
)

type Chat struct {
	gorm.Model
	NextID uint   `json:"next_id"`
	Admin  bool   `json:"admin"`
	Data   string `json:"data"`
}

type ResultDatabase struct {
	Err  error
	Chat []Chat
}
