package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

func check(input core.Memo) error {
	if input.GroupID == 0 {
		return fmt.Errorf("GroupID is wrong... ")
	}
	if !(1 <= input.Type && input.Type <= 3) {
		return fmt.Errorf("type is wrong... ")
	}
	if input.Message == "" {
		return fmt.Errorf("message is wrong... ")
	}

	return nil
}
