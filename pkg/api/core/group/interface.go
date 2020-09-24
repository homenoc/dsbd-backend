package group

const (
	ID           = 0
	OrgJa        = 1
	Org          = 2
	Email        = 3
	UpdateID     = 100
	UpdateOrg    = 101
	UpdateStatus = 102
	UpdateTechID = 103
	UpdateInfo   = 104
)

type Group struct {
	ID        int    `json:"id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	OrgJa     int    `json:"org_ja"`
	Org       string `json:"org"`
	Status    string `json:"status"`
	TechID    string `json:"tech_id"`
	PostCode  int    `json:"postcode"`
	AddressJa int    `json:"address_ja"`
	Address   int    `json:"address"`
	Mail      string `json:"mail"`
	Phone     string `json:"phone"`
}

type Result struct {
	Status    bool    `json:"status"`
	Error     string  `json:"error"`
	GroupData []Group `json:"data"`
}
