package payment

import "github.com/homenoc/dsbd-backend/pkg/api/core"

const (
	ID              = 0
	PaymentIntentID = 1
	UpdatePaid      = 140
	UpdateAll       = 150
)

type Input struct {
	Plan string `json:"plan"` // yearly, monthly
}

type ChangeCardPaymentInit struct {
	PaymentMethodID string `json:"payment_method_id"`
}

type ResultByAdmin struct {
	Payment []core.Payment `json:"payment"`
}

type ResultByUser struct {
	ClientSecret string `json:"client_secret"`
}
