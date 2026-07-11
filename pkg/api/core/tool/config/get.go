package config

import "fmt"

// Service/connection type definitions moved to the code registry in package
// core (core.GetServiceType / core.GetConnectionType); they are no longer in
// config because each capability flag drives Go validation and a frontend form.

func GetMailTemplate(id string) (*MailTemplate, error) {
	for _, mail := range Conf.Template.Mail {
		if mail.ID == id {
			return &mail, nil
		}
	}

	// Return a non-nil empty template on not-found: every caller ignores the
	// error and dereferences the result (e.g. mailTemplate.Message), so a nil
	// here panics (500) — as hit during account creation when no "signature"
	// template is configured.
	return &MailTemplate{}, fmt.Errorf("mail template is not found")
}

func GetMembershipTemplate(plan string) (*MembershipTemplate, error) {
	for _, membership := range Conf.Template.Membership {
		if membership.Plan == plan {
			return &membership, nil
		}
	}

	return nil, fmt.Errorf("mail template is not found")
}
