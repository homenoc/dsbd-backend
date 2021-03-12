package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

func replace(input, replace core.NOC) core.NOC {

	//Name
	if input.Name != "" {
		replace.Name = input.Name
	}
	//Location
	if input.Location != "" {
		replace.Location = input.Location
	}
	//Bandwidth
	if input.Bandwidth != "" {
		replace.Bandwidth = input.Bandwidth
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
