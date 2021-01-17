package v0

import (
	router "github.com/homenoc/dsbd-backend/pkg/api/core/noc/router"
)

func replace(input, replace router.Router) router.Router {

	//HostName
	if input.HostName != "" {
		replace.HostName = input.HostName
	}
	//Address
	if input.Address != "" {
		replace.Address = input.Address
	}

	// uint boolean
	//HostName
	if input.NOC != 0 {
		replace.NOC = input.NOC
	}
	//Enable
	if input.Enable != replace.Enable {
		replace.Enable = input.Enable
	}

	return replace
}
