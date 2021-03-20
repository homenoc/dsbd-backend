package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

func replace(input, replace core.BGPRouter) core.BGPRouter {

	//HostName
	if input.HostName != "" {
		replace.HostName = input.HostName
	}
	//Address
	if input.Address != "" {
		replace.Address = input.Address
	}

	// uint boolean
	//NOC
	if input.NOCID != 0 {
		replace.NOCID = input.NOCID
	}
	//Enable
	if input.Enable != replace.Enable {
		replace.Enable = input.Enable
	}

	return replace
}
