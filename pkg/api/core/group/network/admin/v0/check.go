package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/admin"
)

func check(input admin.Admin) error {
	// check
	if input.UserID <= 0 {
		return fmt.Errorf("failed data: user id")
	}
	if input.NetworkID <= 0 {
		return fmt.Errorf("failed data: network id")
	}
	return nil
}
