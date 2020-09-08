package auth

type User struct {
	ID        int
	CreatedAt int
	UpdatedAt int
	GID       int
	Name      string
	Pass      string
	Mail      string
	Level     int
	Status    int
	IsVerify  int
	MailToken string
}

type Group struct {
	ID        int
	CreatedAt int
	UpdatedAt int
	OrgJa     string
	Org       string
	Status    int
	TechID    string
	PostCode  string
	AddressJa string
	Address   string
	Mail      string
	Phone     string
}

type Auth struct {
	User  *User
	Group *Group
}

type ResultAuth struct {
	Status bool
	Error  string
}
