package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

func replace(input, replace core.TunnelEndPointRouter) core.TunnelEndPointRouter {

	//HostName
	if input.HostName != "" {
		replace.HostName = input.HostName
	}

	//Comment
	if input.Comment != "" {
		replace.Comment = input.Comment
	}

	// uint boolean
	//NOCID
	if input.NOCID != 0 {
		replace.NOCID = input.NOCID
	}
	//Capacity
	if input.Capacity != 0 {
		replace.Capacity = input.Capacity
	}
	//Enable
	if input.Enable != replace.Enable {
		replace.Enable = input.Enable
	}

	return replace
}
