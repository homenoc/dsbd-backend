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
	if input.Org == "" {
		return fmt.Errorf("no data: org")
	}
	if input.Contract == "" {
		return fmt.Errorf("no data: contract")
	}
	return nil
}
