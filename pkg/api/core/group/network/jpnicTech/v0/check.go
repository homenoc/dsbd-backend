package v0

import (
	"fmt"
	jpnic "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicTech"
)

func check(input jpnic.JpnicTech) error {
	// check
	if input.UserID <= 0 {
		return fmt.Errorf("failed data: user id")
	}
	if input.NetworkID <= 0 {
		return fmt.Errorf("failed data: network id")
	}
	return nil
}
