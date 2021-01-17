package v0

import (
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
)

func replaceNetwork(replace, input network.Network) network.Network {
	//Org
	if input.Org != "" {
		replace.Org = input.Org
	}

	//Org (English)
	if input.OrgEn != "" {
		replace.OrgEn = input.OrgEn
	}

	//Postcode
	if input.Postcode != "" {
		replace.Postcode = input.Postcode
	}

	//Address
	if input.Address != "" {
		replace.Address = input.Address
	}

	//Address(English)
	if input.AddressEn != "" {
		replace.AddressEn = input.AddressEn
	}

	//Route(V4)
	if input.RouteV4 != "" {
		replace.RouteV4 = input.RouteV4
	}

	//Route(V4)
	if input.RouteV6 != "" {
		replace.RouteV6 = input.RouteV6
	}

	//PI
	if input.PI != replace.PI {
		replace.PI = input.PI
	}

	//Lock
	if input.Lock != replace.Lock {
		replace.Lock = input.Lock
	}

	//ASN
	if input.ASN != "" {
		replace.ASN = input.ASN
	}

	//V4
	if input.V4 != "" {
		replace.V4 = input.V4
	}

	//V6
	if input.V6 != "" {
		replace.V6 = input.V6
	}

	//V4Name
	if input.V4Name != "" {
		replace.V4Name = input.V4Name
	}

	//V6Name
	if input.V6Name != "" {
		replace.V6Name = input.V6Name
	}
	//Date
	if input.Date != "" {
		replace.Date = input.Date
	}

	//Plan
	if input.Plan != "" {
		replace.Plan = input.Plan
	}

	return replace
}

func replaceAdminNetwork(replace, input network.Network) network.Network {
	//Org
	if input.Org != "" {
		replace.Org = input.Org
	}

	//Org (English)
	if input.OrgEn != "" {
		replace.OrgEn = input.OrgEn
	}

	//Postcode
	if input.Postcode != "" {
		replace.Postcode = input.Postcode
	}

	//Address
	if input.Address != "" {
		replace.Address = input.Address
	}

	//Address(English)
	if input.AddressEn != "" {
		replace.AddressEn = input.AddressEn
	}

	//Route(V4)
	if input.RouteV4 != "" {
		replace.RouteV4 = input.RouteV4
	}

	//Route(V4)
	if input.RouteV6 != "" {
		replace.RouteV6 = input.RouteV6
	}

	//ASN
	if input.ASN != "" {
		replace.ASN = input.ASN
	}

	//V4
	if input.V4 != "" {
		replace.V4 = input.V4
	}

	//V6
	if input.V6 != "" {
		replace.V6 = input.V6
	}

	//V4Name
	if input.V4Name != "" {
		replace.V4Name = input.V4Name
	}

	//V6Name
	if input.V6Name != "" {
		replace.V6Name = input.V6Name
	}
	//Date
	if input.Date != "" {
		replace.Date = input.Date
	}

	//Plan
	if input.Plan != "" {
		replace.Plan = input.Plan
	}

	// bool
	//Lock
	if input.Lock != replace.Lock {
		replace.Lock = input.Lock
	}

	//PI
	if input.PI != replace.PI {
		replace.PI = input.PI
	}

	//Open
	if input.Open != replace.Open {
		replace.Open = input.Open
	}

	return replace
}
