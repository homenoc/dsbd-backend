package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
)

func check(input connection.Input) error {

	if input.Address == "" {
		return fmt.Errorf("error: address is invalid")
	}

	return nil
}
