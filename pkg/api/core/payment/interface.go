package payment

type Input struct {
	Plan string `json:"plan"` // yearly, monthly
}

type ChangeCardPaymentInit struct {
	PaymentMethodID string `json:"payment_method_id"`
}

type ResultByUser struct {
	ClientSecret string `json:"client_secret"`
}
