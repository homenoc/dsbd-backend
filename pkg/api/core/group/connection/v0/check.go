package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
)

func check(input connection.Input) error {

	if input.Prefectures > 48 {
		return fmt.Errorf("error: prefectures is invalid...")
	}

	return nil
}
