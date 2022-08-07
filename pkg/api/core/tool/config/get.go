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
