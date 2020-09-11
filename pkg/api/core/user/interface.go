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

type InitToken struct {
	ID        int    `json:"id"`
	CreatedAt int    `json:"created_at"`
	ExpiredAt int    `json:"expired_at"`
	DeletedAt int    `json:"deleted_at"`
	Mail      string `json:"mail"`
	Token1    string `json:"token1"`
	Token2    string `json:"token2"`
	IP        string `json:"ip"`
}

type Result struct {
	Status   bool   `json:"status"`
	Error    string `json:"error"`
	UserData []User `json:"data"`
}

type InitTokenResult struct {
	Status    bool        `json:"status"`
	Error     string      `json:"error"`
	TokenData []InitToken `json:"token"`
}
