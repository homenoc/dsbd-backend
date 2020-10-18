package v0

import (
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
)

func replaceNetwork(serverData, input network.Network) network.Network {
	input.ID = serverData.ID

	//Org
	if input.Org == "" {
		input.Org = serverData.Org
	}

	//Org (English)
	if input.OrgEn == "" {
		input.OrgEn = serverData.OrgEn
	}

	//Postcode
	if input.Postcode == "" {
		input.Postcode = serverData.Postcode
	}

	//Address
	if input.Address == "" {
		input.Address = serverData.Address
	}

	//Address(English)
	if input.AddressEn == "" {
		input.AddressEn = serverData.AddressEn
	}

	//Route
	if input.Route == "" {
		input.Route = serverData.Route
	}

	//PI
	input.PI = serverData.PI

	//Lock
	input.Lock = serverData.Lock

	//ASN
	if input.ASN == "" {
		input.ASN = serverData.ASN
	}

	//V4
	if input.V4 == "" {
		input.V4 = serverData.V4
	}

	//V6
	if input.V6 == "" {
		input.V6 = serverData.V6
	}

	//V4Name
	if input.V4Name == "" {
		input.V4Name = serverData.V4Name
	}

	//V6Name
	if input.V6Name == "" {
		input.V6Name = serverData.V6Name
	}
	//Date
	if input.Date == "" {
		input.Date = serverData.Date
	}

	//Plan
	if input.Plan == "" {
		input.Plan = serverData.Plan
	}

	return input
}
