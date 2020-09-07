package user

const (
	ID        = 0
	GID       = 1
	Name      = 2
	Email     = 3
	MailToken = 4
)

type User struct {
	ID         int    `json:"id"`
	CreatedAt  int    `json:"create_time"`
	UpdatedAt  int    `json:"update_at"`
	GID        int    `json:"gid"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Pass       string `json:"pass"`
	Status     int    `json:"status"`
	Level      int    `json:"level"`
	MailVerify bool   `json:"mail_verify"`
	MailToken  string `json:"mail_token"`
}

type Result struct {
	Status   string `json:"result"`
	Data     string `json:"token"`
	UserData []User `json:"data"`
}
