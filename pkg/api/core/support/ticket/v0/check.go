package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support"
)

func check(input support.FirstInput) error {
	if input.Title == "" {
		return fmt.Errorf("no data: title")
	}

	if input.Data == "" {
		return fmt.Errorf("no data: data")
	}

	return nil
}

func checkByAdmin(input support.FirstInput) error {
	if input.Title == "" {
		return fmt.Errorf("no data: title")
	}
	if input.Data == "" {
		return fmt.Errorf("no data: data")
	}

	return nil
}
