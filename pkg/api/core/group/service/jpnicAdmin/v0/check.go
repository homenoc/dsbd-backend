package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

func check(input core.JPNICAdmin) error {
	// check
	if input.Org == "" {
		return fmt.Errorf("failed data: org")
	}
	if input.OrgEn == "" {
		return fmt.Errorf("failed data: org(english)")
	}
	if input.PostCode == "" {
		return fmt.Errorf("failed data: postcode")
	}
	if input.Address == "" {
		return fmt.Errorf("failed data: address")
	}
	if input.AddressEn == "" {
		return fmt.Errorf("failed data: address(english)")
	}
	if input.Tel == "" {
		return fmt.Errorf("failed data: tel")
	}
	return nil
}
