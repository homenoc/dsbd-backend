package service

import "github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"

type Result struct {
	Network    []config.Network    `json:"network"`
	Connection []config.Connection `json:"connection"`
}
