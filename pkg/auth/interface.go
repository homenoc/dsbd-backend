package auth

type User struct {
	ID        int
	CreatedAt int
	UpdatedAt int
	Name      string
	GID       int
	Pass      string
	Mail      string
	Level     int
	Status    int
	IsVerify  int
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
