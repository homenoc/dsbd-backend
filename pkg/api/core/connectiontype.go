package core

import "fmt"

// ConnectionType is a connection kind (EtherIP / GRE / IP-IP / Cross Connect /
// IXP / etc). Registry, not config — see ServiceType for the rationale.
//
// NeedInternet is the only flag consumed by Go (connection validation); the
// rest (NeedComment/NeedCrossConnect/IsL2/IsL3) are used only by the frontend
// form and are mirrored in packages/shared/constants.
type ConnectionType struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	Comment          string `json:"comment"`
	NeedInternet     bool   `json:"need_internet"`
	NeedComment      bool   `json:"need_comment"`
	NeedCrossConnect bool   `json:"need_cross_connect"`
	IsL2             bool   `json:"is_l2"`
	IsL3             bool   `json:"is_l3"`
}

// connectionTypes seeds from the former config.json template.connection values.
var connectionTypes = []ConnectionType{
	{Name: "Ethernet over IPトンネル", Type: "EIP", NeedInternet: true, IsL2: true},
	{Name: "GREトンネル", Type: "GRE", NeedInternet: true, IsL3: true},
	{Name: "IP-IPトンネル", Type: "IPT", NeedInternet: true, IsL3: true},
	{Name: "直接接続", Type: "CC0", Comment: "データセンター構内回線、キャリア専用線上などで直接接続",
		NeedComment: true, NeedCrossConnect: true, IsL2: true, IsL3: true},
	{Name: "IX接続", Type: "IXP", Comment: "インターネットエクスチェンジ経由で接続",
		NeedComment: true, NeedCrossConnect: true, IsL2: true, IsL3: true},
	{Name: "その他", Type: "ETC", Comment: "その他の接続方式", NeedComment: true},
}

// ConnectionTypes returns the full registry.
func ConnectionTypes() []ConnectionType { return connectionTypes }

// GetConnectionType looks up a connection type by its code (e.g. "EIP").
func GetConnectionType(t string) (ConnectionType, error) {
	for _, ct := range connectionTypes {
		if ct.Type == t {
			return ct, nil
		}
	}
	return ConnectionType{}, fmt.Errorf("connection type not found: %s", t)
}
