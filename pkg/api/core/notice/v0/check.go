package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
)

func check(input notice.Notice) error {
	// check
	if input.Title == "" {
		return fmt.Errorf("no data: title")
	}
	if input.Data == "" {
		return fmt.Errorf("no data: data")
	}
	return nil
}
