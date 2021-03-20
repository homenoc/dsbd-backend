package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
)

func replace(baseData *core.Group, input group.Input) *core.Group {
	if input.Org != "" {
		baseData.Org = input.Org
	}
	if input.OrgEn != "" {
		baseData.OrgEn = input.OrgEn
	}
	if input.PostCode != "" {
		baseData.PostCode = input.PostCode
	}
	if input.Address != "" {
		baseData.Address = input.Address
	}
	if input.AddressEn != "" {
		baseData.AddressEn = input.AddressEn
	}
	if input.Tel != "" {
		baseData.Tel = input.Tel
	}
	if input.Country != "" {
		baseData.Country = input.Country
	}

	return baseData
}
