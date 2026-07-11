package core

import "fmt"

// ServiceType is a network service kind (L2 / L3 Static / L3 BGP / Transit).
// This registry is the source of truth for service types and their capability
// flags — moved out of config.json into code because every flag drives Go
// validation and a frontend form section, so adding a type is a code change by
// definition. The frontend mirrors this list in packages/shared/constants.
//
// NetBox/provisioning metadata (feature a/b) will attach here as new fields.
type ServiceType struct {
	Hidden       bool   `json:"hidden"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Comment      string `json:"comment"`
	NeedJPNIC    bool   `json:"need_jpnic"`
	NeedGlobalAS bool   `json:"need_global_as"`
	NeedComment  bool   `json:"need_comment"`
	NeedRoute    bool   `json:"need_route"`
	NeedBGP      bool   `json:"need_bgp"`
}

// serviceTypes seeds from the former config.json template.service values.
var serviceTypes = []ServiceType{
	{Name: "L2", Type: "2000", Comment: "L2接続", NeedJPNIC: true},
	{Name: "L3 Static", Type: "3S00", Comment: "Staticにて接続", NeedJPNIC: true, NeedRoute: true},
	{Name: "L3 BGP", Type: "3B00", Comment: "BGPにて接続", NeedJPNIC: true, NeedRoute: true, NeedBGP: true},
	{Name: "トランジット提供", Type: "IP3B", Comment: "グローバルASのユーザに対するトランジット提供(L3 BGP)",
		NeedGlobalAS: true, NeedRoute: true, NeedBGP: true},
}

// ServiceTypes returns the full registry.
func ServiceTypes() []ServiceType { return serviceTypes }

// GetServiceType looks up a service type by its code (e.g. "3B00").
func GetServiceType(t string) (ServiceType, error) {
	for _, st := range serviceTypes {
		if st.Type == t {
			return st, nil
		}
	}
	return ServiceType{}, fmt.Errorf("service type not found: %s", t)
}
