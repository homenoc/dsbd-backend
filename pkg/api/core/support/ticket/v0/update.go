package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
)

func updateAdminTicket(input, replace ticket.Ticket) (ticket.Ticket, error) {

	//Title
	if input.Title != "" {
		replace.Title = input.Title
	}
	// uint boolean
	//Solved
	if input.Solved != replace.Solved {
		replace.Solved = input.Solved
	}
	//UserID
	if input.UserID != replace.UserID {
		replace.UserID = input.UserID
	}
	//GroupID
	if input.GroupID != replace.GroupID {
		replace.GroupID = input.GroupID
	}

	return replace, nil
}