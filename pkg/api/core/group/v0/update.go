package v0

import "github.com/homenoc/dsbd-backend/pkg/api/core/group"

func updateAdminUser(input, replace group.Group) (group.Group, error) {

	//Org
	if input.Org != "" {
		replace.Status = input.Status
	}

	// uint boolean
	//Lock
	if input.Lock != replace.Lock {
		replace.Lock = input.Lock
	}
	//Status
	if input.Status != replace.Status {
		replace.Status = input.Status
	}

	return replace, nil
}
