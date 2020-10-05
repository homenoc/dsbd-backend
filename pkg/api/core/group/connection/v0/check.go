package v0

import (
	"fmt"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
)

func check(input connection.Connection) error {
	// check
	if input.NTT == "" {
		return fmt.Errorf("no data: NTT")
	}
	if input.Service == "" {
		return fmt.Errorf("no data: service")
	}
	if input.NOC == "" {
		return fmt.Errorf("no data: noc")
	}
	if input.TermIP == "" {
		return fmt.Errorf("no data: term ip")
	}
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	if input.Org == "" {
		return fmt.Errorf("no data: org")
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
	return nil
}
