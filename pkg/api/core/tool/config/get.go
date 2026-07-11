package config

import "fmt"

func GetMailTemplate(id string) (*MailTemplate, error) {
	for _, mail := range Conf.Template.Mail {
		if mail.ID == id {
			return &mail, nil
		}
	}

	// Return a non-nil empty template: callers append .Message to mail bodies
	// inside long-lived loops (the support WebSocket handlers), where a nil
	// deref panics and tears down the connection.
	return &MailTemplate{}, fmt.Errorf("mail template is not found")
}

func GetMembershipTemplate(plan string) (*MembershipTemplate, error) {
	for _, membership := range Conf.Template.Membership {
		if membership.Plan == plan {
			return &membership, nil
		}
	}

	// Non-nil for the same defensive reason as GetMailTemplate.
	return &MembershipTemplate{}, fmt.Errorf("membership template is not found")
}
