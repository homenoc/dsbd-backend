package gatewayIP

import (
	"github.com/jinzhu/gorm"
)

const (
	ID        = 0
	NOC       = 1
	Address   = 2
	Enable    = 3
	UpdateAll = 110
)

type GatewayIP struct {
	gorm.Model
	GatewayID uint   `json:"gateway_id"`
	IP        string `json:"ip"`
	Comment   string `json:"comment"`
	Enable    bool   `json:"enable"`
}

type Result struct {
	GatewayIP []GatewayIP `json:"gateway_ip"`
}

type ResultDatabase struct {
	Err       error
	GatewayIP []GatewayIP
}
