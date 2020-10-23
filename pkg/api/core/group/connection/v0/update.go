package v0

import (
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
)

func updateAdminConnection(input, replace connection.Connection) (connection.Connection, error) {

	//Service ID
	if input.ServiceID != "" {
		replace.ServiceID = input.ServiceID
	}
	//Service
	if input.Service != "" {
		replace.Service = input.Service
	}

	//NOC
	if input.NOC != "" {
		replace.NOC = input.NOC
	}
	//NOC IP
	if input.NOCIP != "" {
		replace.NOCIP = input.NOCIP
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
	// Open
	if input.Open != replace.Open {
		replace.Open = input.Open
	}
	// Fee
	if input.Fee != replace.Fee {
		replace.Fee = input.Fee
	}

	return replace, nil
}
