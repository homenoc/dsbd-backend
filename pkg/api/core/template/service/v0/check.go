package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

func check(input core.ServiceTemplate) error {
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}

	if input.Comment == "" {
		return fmt.Errorf("no data: comment")
	}

	return nil
}
