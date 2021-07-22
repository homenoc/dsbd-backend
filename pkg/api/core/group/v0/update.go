package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

func updateAdminGroup(input, replace core.Group) (core.Group, error) {

	//Org
	if input.Org != "" {
		replace.Org = input.Org
	}
	// Org(English)
	if input.OrgEn != "" {
		replace.OrgEn = input.OrgEn
	}
	// PostCode
	if input.PostCode != "" {
		replace.PostCode = input.PostCode
	}
	// Address
	if input.Address != "" {
		replace.Address = input.Address
	}
	// Address(English)
	if input.AddressEn != "" {
		replace.AddressEn = input.AddressEn
	}
	// Tel
	if input.Tel != "" {
		replace.Tel = input.Tel
	}
	// Country
	if input.Country != "" {
		replace.Country = input.Country
	}

	// Pass
	if input.Pass != replace.Pass {
		replace.Pass = input.Pass
	}
	// ExpiredStatus
	if input.ExpiredStatus != replace.ExpiredStatus {
		replace.ExpiredStatus = input.ExpiredStatus
	}

	return replace, nil
}
