package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

func replaceService(replace, input core.Service) core.Service {
	//Org
	if input.Org != "" {
		replace.Org = input.Org
	}

	//Org (English)
	if input.OrgEn != "" {
		replace.OrgEn = input.OrgEn
	}

	//PostCode
	if input.PostCode != "" {
		replace.PostCode = input.PostCode
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

	//Route(V6)
	if input.RouteV6 != "" {
		replace.RouteV6 = input.RouteV6
	}

	//Open
	if input.Pass != replace.Pass {
		replace.Pass = input.Pass
	}

	//Lock
	if input.Lock != replace.Lock {
		replace.Lock = input.Lock
	}

	////V4
	//if input.V4 != "" {
	//	replace.V4 = input.V4
	//}
	//
	////V6
	//if input.V6 != "" {
	//	replace.V6 = input.V6
	//}
	//
	////V4Name
	//if input.V4Name != "" {
	//	replace.V4Name = input.V4Name
	//}
	//
	////V6Name
	//if input.V6Name != "" {
	//	replace.V6Name = input.V6Name
	//}
	////Date
	//if input.Date != "" {
	//	replace.Date = input.Date
	//}

	return replace
}
