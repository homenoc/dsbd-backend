package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
)

func check(input jpnicAdmin.JpnicAdmin) error {
	// check
	if input.UserID <= 0 {
		return fmt.Errorf("failed data: user id")
	}
	if input.NetworkID <= 0 {
		return fmt.Errorf("failed data: network id")
	}
	return nil
}
