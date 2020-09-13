package token

const (
	ID                      = 0
	UserToken               = 10
	UserTokenAndAccessToken = 11
	AddToken                = 100
	UpdateToken             = 101
)

type Token struct {
	ID          int    `json:"id"`
	CreatedAt   int    `json:"created_at"`
	UpdatedAt   int    `json:"updated_at"`
	ExpiredAt   int    `json:"expired_at"`
	DeletedAt   int    `json:"deleted_at"`
	UID         int    `json:"uid"`
	Status      int    `json:"status"` //0: before login 1: after login
	UserToken   string `json:"user_token"`
	TmpToken    string `json:"tmp_token"`
	AccessToken string `json:"access_token"`
	Debug       string `json:"debug"`
}

type Result struct {
	Status bool    `json:"status"`
	Error  string  `json:"error"`
	Token  []Token `json:"token"`
}
