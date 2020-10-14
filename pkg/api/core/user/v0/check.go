package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
)

func check(input user.User) error {
	// check
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	if input.Org == "" {
		return fmt.Errorf("no data: position")
	}
	if input.PostCode == "" {
		return fmt.Errorf("no data: postcode")
	}
	if input.Address == "" {
		return fmt.Errorf("no data: address")
	}
	if input.Phone == "" {
		return fmt.Errorf("no data: phone")
	}
	if input.Country == "" {
		return fmt.Errorf("no data: country")
	}
	return nil
}
