package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

func updateAdminConnection(input, replace core.Connection) core.Connection {
	//Term IP
	if input.TermIP != "" {
		replace.TermIP = input.TermIP
	}

	//LinkV4Our
	if input.LinkV4Our != "" {
		replace.LinkV4Our = input.LinkV4Our
	}
	//LinkV4Your
	if input.LinkV4Your != "" {
		replace.LinkV4Your = input.LinkV4Your
	}
	//LinkV6Our
	if input.LinkV6Our != "" {
		replace.LinkV6Our = input.LinkV6Our
	}
	//LinkV6Your
	if input.LinkV6Your != "" {
		replace.LinkV6Your = input.LinkV6Your
	}

	//Fee
	if input.Fee != "" {
		replace.Fee = input.Fee
	}

	// uint boolean
	replace.ConnectionNumber = input.ConnectionNumber

	// Open
	replace.Open = input.Open

	// Monitor
	replace.Monitor = input.Monitor

	//ServiceType
	replace.ConnectionTemplateID = input.ConnectionTemplateID

	//NTT
	replace.NTTTemplateID = input.NTTTemplateID

	//NOC
	replace.NOCID = input.NOCID

	return replace
}
