package v0

import (
	"fmt"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
)

func check(input network.Network) error {
	// check
	if input.Route == "" {
		return fmt.Errorf("no data: route")
	}
	if input.Date == "" {
		return fmt.Errorf("no data: date")
	}
	return nil
}
