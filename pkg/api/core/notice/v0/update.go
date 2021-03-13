package v0

import "github.com/homenoc/dsbd-backend/pkg/api/core"

func updateAdminUser(input, replace core.Notice) core.Notice {

	//Title
	if input.Title != "" {
		replace.Title = input.Title
	}
	//Data
	if input.Data != "" {
		replace.Data = input.Data
	}

	// uint boolean
	//UserID
	if input.UserID != replace.UserID {
		replace.UserID = input.UserID
	}
	//GroupID
	if input.GroupID != replace.GroupID {
		replace.GroupID = input.GroupID
	}
	//StartTime
	if input.StartTime != replace.StartTime {
		replace.StartTime = input.StartTime
	}
	//EndTime
	if input.EndingTime != replace.EndingTime {
		replace.EndingTime = input.EndingTime
	}
	//Everyone
	if input.Everyone != replace.Everyone {
		replace.Everyone = input.Everyone
	}
	//Important
	if input.Important != replace.Important {
		replace.Important = input.Important
	}
	//Fault
	if input.Fault != replace.Fault {
		replace.Fault = input.Fault
	}
	//Info
	if input.Info != replace.Info {
		replace.Info = input.Info
	}

	return replace
}
