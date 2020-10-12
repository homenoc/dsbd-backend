package v0

import (
	"fmt"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
)

func check(input network.Network) error {
	// check
	if !(1 <= input.Type && input.Type <= 2 || input.Type == 5) {
		return fmt.Errorf("error: type value")
	}
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	if input.IP == "" {
		return fmt.Errorf("no data: ip")
	}
	if input.Route == "" {
		return fmt.Errorf("no data: route")
	}
	if input.Date == "" {
		return fmt.Errorf("no data: date")
	}
	if input.Plan == "" {
		return fmt.Errorf("no data: plan")
	}
	return nil
}

func checkUpdate(input network.NetworkUser) error {
	// check
	if !(1 <= input.Type && input.Type <= 2 || input.Type == 5) {
		return fmt.Errorf("error: type value")
	}
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	if input.IP == "" {
		return fmt.Errorf("no data: ip")
	}
	if input.Route == "" {
		return fmt.Errorf("no data: route")
	}
	if input.Date == "" {
		return fmt.Errorf("no data: date")
	}
	if input.Plan == "" {
		return fmt.Errorf("no data: plan")
	}
	if len(input.OperationID) == 0 {
		return fmt.Errorf("no data: operationID")
	}
	if len(input.TechID) == 0 {
		return fmt.Errorf("no data: techID")
	}
	return nil
}
