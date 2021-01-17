package v0

import (
	router "github.com/homenoc/dsbd-backend/pkg/api/core/noc/router"
)

func replace(input, replace router.Router) router.Router {

	//Host
	if input.Host != "" {
		replace.Host = input.Host
	}
	//Address
	if input.Address != "" {
		replace.Address = input.Address
	}

	// uint boolean
	//Host
	if input.NOC != 0 {
		replace.NOC = input.NOC
	}
	//Enable
	if input.Enable != replace.Enable {
		replace.Enable = input.Enable
	}

	return replace
}
