package v0

import "github.com/homenoc/dsbd-backend/pkg/api/core"

func replace(input, replace core.TunnelEndPointRouterIP) core.TunnelEndPointRouterIP {

	//IP
	if input.IP != "" {
		replace.IP = input.IP
	}
	//Comment
	if input.Comment != "" {
		replace.Comment = input.Comment
	}

	// uint boolean
	//Enable
	if input.Enable != replace.Enable {
		replace.Enable = input.Enable
	}

	return replace
}
