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
	if input.UserId == 0 {
		return fmt.Errorf("no data: userID")
	}
	return nil
}
