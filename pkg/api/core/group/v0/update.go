package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
)

func updateAdminGroup(input, replace group.Group) (group.Group, error) {

	//Org
	if input.Org != "" {
		replace.Org = input.Org
	}

	// uint boolean
	// Lock
	if input.Lock != replace.Lock {
		replace.Lock = input.Lock
	}
	// Pass
	if input.Pass != replace.Pass {
		replace.Pass = input.Pass
	}
	// Status
	if input.Status != replace.Status {
		replace.Status = input.Status
	}
	// ExpiredStatus
	if input.ExpiredStatus != replace.ExpiredStatus {
		replace.ExpiredStatus = input.ExpiredStatus
	}

	return replace, nil
}
