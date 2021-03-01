package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
)

func updateAdminConnection(input, replace connection.Connection) connection.Connection {
	//ServiceType
	if input.ConnectionType != "" {
		replace.ConnectionType = input.ConnectionType
	}

	//NTT
	if input.NTT != "" {
		replace.NTT = input.NTT
	}

	//NOC
	if input.NOC != "" {
		replace.NOC = input.NOC
	}

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
	replace.NetworkID = input.NetworkID
	replace.GroupID = input.GroupID
	replace.UserID = input.UserID
	replace.GatewayIPID = input.GatewayIPID
	replace.RouterID = input.RouterID
	replace.ConnectionNumber = input.ConnectionNumber

	// Open
	replace.Open = input.Open

	// Monitor
	replace.Monitor = input.Monitor

	return replace
}
