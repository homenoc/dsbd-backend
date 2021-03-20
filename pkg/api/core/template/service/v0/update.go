package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

func update(input, replace core.ServiceTemplate) core.ServiceTemplate {

	//Name
	if input.Name != "" {
		replace.Name = input.Name
	}

	//Comment
	if input.Comment != "" {
		replace.Comment = input.Comment
	}

	// uint boolean
	//Hidden
	if input.Hidden != replace.Hidden {
		replace.Hidden = input.Hidden
	}

	return replace
}
