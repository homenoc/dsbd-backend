package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"strings"
)

func check(input user.Input) error {
	// check
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}

	if !strings.Contains(input.Email, "@") {
		return fmt.Errorf("no data: email address")
	}
	if input.Name == "" || input.NameEn == "" {
		return fmt.Errorf("no data: name")
	}

	return nil
}
