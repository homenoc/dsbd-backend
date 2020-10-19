package v0

import (
	"fmt"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
)

func check(input network.NetworkInput) error {
	// check
	if input.Route == "" {
		return fmt.Errorf("no data: route")
	}
	if input.Date == "" {
		return fmt.Errorf("no data: date")
	}
	if input.PI {
		if input.ASN == "" {
			return fmt.Errorf("no data: ASN")
		}
	} else {
		if input.Org == "" {
			return fmt.Errorf("no data: Org")
		}
		if input.OrgEn == "" {
			return fmt.Errorf("no data: Org(English)")
		}
		if input.Postcode == "" {
			return fmt.Errorf("no data: postcode")
		}
		if input.Address == "" {
			return fmt.Errorf("no data: Address")
		}
		if input.AddressEn == "" {
			return fmt.Errorf("no data: Address(English)")
		}
		if input.V4 == "" {
			return fmt.Errorf("no data: v4")
		}
		if input.V6 == "" {
			return fmt.Errorf("no data: v6")
		}
		if input.V4Name == "" {
			return fmt.Errorf("no data: v4 Name")
		}
		if input.V6Name == "" {
			return fmt.Errorf("no data: v6 Name")
		}
	}
	return nil
}
