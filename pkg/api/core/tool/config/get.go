package config

import "fmt"

func GetMailTemplate(id string) (*MailTemplate, error) {
	for _, mail := range Conf.Template.Mail {
		if mail.ID == id {
			return &mail, nil
		}
	}

	return nil, fmt.Errorf("mail template is not found")
}

func GetMembershipTemplate(plan string) (*MembershipTemplate, error) {
	for _, membership := range Conf.Template.Membership {
		if membership.Plan == plan {
			return &membership, nil
		}
	}

	return nil, fmt.Errorf("mail template is not found")
}
