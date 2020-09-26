package user

const (
	ID               = 0
	GID              = 1
	Name             = 2
	Email            = 3
	MailToken        = 4
	UpdateVerifyMail = 100
	UpdateGID        = 101
	UpdateName       = 102
	UpdateMail       = 103
	UpdatePass       = 104
	UpdateStatus     = 105
	UpdateLevel      = 106
)

type User struct {
	ID         int    `json:"id"`
	CreatedAt  int    `json:"created_at"`
	UpdatedAt  int    `json:"updated_at"`
	GID        int    `json:"gid"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Pass       string `json:"pass"`
	Status     int    `json:"status"`
	Level      int    `json:"level"`
	MailVerify int    `json:"mail_verify"`
	MailToken  string `json:"mail_token"`
}

type Result struct {
	Status   bool   `json:"status"`
	Error    string `json:"error"`
	UserData []User `json:"data"`
}
