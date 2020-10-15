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
	return nil
}
