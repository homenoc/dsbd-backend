package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/gatewayIP"
)

func replace(input, replace gatewayIP.GatewayIP) gatewayIP.GatewayIP {

	//IP
	if input.IP != "" {
		replace.IP = input.IP
	}
	//Comment
	if input.Comment != "" {
		replace.Comment = input.Comment
	}

	// uint boolean
	//GatewayID
	if input.GatewayID != 0 {
		replace.GatewayID = input.GatewayID
	}
	//Enable
	if input.Enable != replace.Enable {
		replace.Enable = input.Enable
	}

	return replace
}
