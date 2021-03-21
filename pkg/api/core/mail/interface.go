package mail

type Mail struct {
	ToMail  string `json:"to_mail"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}
