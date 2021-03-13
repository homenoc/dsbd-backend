package info

type Info struct {
	ServiceID  string   `json:"service_id"`
	Service    string   `json:"service"`
	ASN        uint     `json:"asn"`
	V4         []string `json:"v4"`
	V6         []string `json:"v6"`
	UserID     uint     `json:"user_id"`
	NOC        string   `json:"noc"`
	NOCIP      string   `json:"noc_ip"`
	TermIP     string   `json:"term_ip"`
	LinkV4Our  string   `json:"link_v4_our"`
	LinkV4Your string   `json:"link_v4_your"`
	LinkV6Our  string   `json:"link_v6_our"`
	LinkV6Your string   `json:"link_v6_your"`
	Fee        string   `json:"fee"`
}

type Result struct {
	Info []Info `json:"info"`
}
