package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

func check(input core.IPv4Template) error {
	if input.Title == "" {
		return fmt.Errorf("no data: title")
	}

	if input.Subnet == "" {
		return fmt.Errorf("no data: subnet")
	}

	if input.Quantity == 0 {
		return fmt.Errorf("no data: quanitity")
	}

	if input.Hide == nil {
		return fmt.Errorf("no data: hide")
	}

	return nil
}
