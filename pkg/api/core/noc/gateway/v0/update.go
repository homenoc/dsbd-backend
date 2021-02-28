package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/gateway"
)

func replace(input, replace gateway.Gateway) gateway.Gateway {

	//HostName
	if input.HostName != "" {
		replace.HostName = input.HostName
	}
	//V4
	if input.V4 != "" {
		replace.V4 = input.V4
	}
	//V6
	if input.V6 != "" {
		replace.V6 = input.V6
	}

	// uint boolean
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
