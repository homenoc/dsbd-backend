package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/mail"
	"strings"
)

func check(input mail.Mail) error {
	// check
	if !strings.Contains(input.ToMail, "@") {
		return fmt.Errorf("invalid: email address")
	}

	if input.Subject == "" {
		return fmt.Errorf("no data: subject")
	}

	if input.Content == "" {
		return fmt.Errorf("no data: content")
	}

	return nil
}
