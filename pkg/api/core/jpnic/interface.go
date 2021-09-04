package jpnic

type GetAllInput struct {
	Version uint   `json:"version"`
	Org     string `json:"org"`
}

type ReturnInput struct {
	Version     uint     `json:"version"`
	Address     []string `json:"address"`
	NetworkName string   `json:"network_name"`
	ReturnDate  string   `json:"return_date"`
	NotifyEMail string   `json:"notify_e_mail"`
}
