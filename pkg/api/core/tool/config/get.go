package config

import "fmt"

func GetServiceTemplate(serviceType string) (*ServiceTemplate, error) {
	for _, service := range Conf.Template.Service {
		if service.Type == serviceType {
			return &service, nil
		}
	}

	return nil, fmt.Errorf("service template is not found")
}

func GetConnectionTemplate(connectionType string) (*ConnectionTemplate, error) {
	for _, connection := range Conf.Template.Connection {
		if connection.Type == connectionType {
			return &connection, nil
		}
	}

	return nil, fmt.Errorf("service template is not found")
}

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
