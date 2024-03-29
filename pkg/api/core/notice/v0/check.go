package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
)

func check(input notice.Input) error {
	// check
	if input.Title == "" {
		return fmt.Errorf("no data: title")
	}
	if input.Body == "" {
		return fmt.Errorf("no data: data")
	}
	return nil
}
