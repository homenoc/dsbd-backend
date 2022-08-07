package config

import "fmt"

func CheckIncludeNTTTemplate(data string) error {
	for _, ntt := range Conf.Template.NTT {
		// configにetcが含む場合は全て弾かない
		if ntt == "etc" || ntt == data {
			return nil
		}
	}

	return fmt.Errorf("ntt template is not found")
}

func CheckIncludeV4RouteTemplate(data string) error {
	for _, v4Route := range Conf.Template.V4Route {
		// configにetcが含む場合は全て弾かない
		if v4Route == "etc" || v4Route == data {
			return nil
		}
	}

	return fmt.Errorf("v4route template is not found")
}

func CheckIncludeV6RouteTemplate(data string) error {
	for _, v6Route := range Conf.Template.V6Route {
		// configにetcが含む場合は全て弾かない
		if v6Route == "etc" || v6Route == data {
			return nil
		}
	}

	return fmt.Errorf("v6route template is not found")
}

func CheckIncludeV4Template(data string) error {
	for _, v4 := range Conf.Template.V4Route {
		// configにetcが含む場合は全て弾かない
		if v4 == "etc" || v4 == data {
			return nil
		}
	}

	return fmt.Errorf("v4 template is not found")
}

func CheckIncludeV6Template(data string) error {
	for _, v6 := range Conf.Template.V6Route {
		// configにetcが含む場合は全て弾かない
		if v6 == "etc" || v6 == data {
			return nil
		}
	}

	return fmt.Errorf("v6 template is not found")
}
