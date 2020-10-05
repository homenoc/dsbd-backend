package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
)

func check(input group.Group) error {
	// check
	if input.Question == "" {
		return fmt.Errorf("no data: question")
	}
	if input.Bandwidth == "" {
		return fmt.Errorf("no data: bandwidth")
	}
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
	if input.Mail == "" {
		return fmt.Errorf("no data: mail")
	}
	if input.Phone == "" {
		return fmt.Errorf("no data: phone")
	}
	if input.Country == "" {
		return fmt.Errorf("no data: country")
	}
	if input.Contract > 2 {
		return fmt.Errorf("no data: failed contract value")
	}
	return nil
}
