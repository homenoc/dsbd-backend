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

func CheckIncludePreferredAPTemplate(data string) error {
	for _, preferredAP := range Conf.Template.PreferredAP {
		if preferredAP == data {
			return nil
		}
	}

	return fmt.Errorf("preferredAP template is not found")
}

func CheckIncludeIXTemplate(data string) error {
	if data == "" {
		return fmt.Errorf("IX is required for IXP connection")
	}
	for _, ix := range Conf.Template.IX {
		if ix.Name == data {
			return nil
		}
	}

	return fmt.Errorf("IX template is not found: %s", data)
}

func CheckIXPeerType(data string) error {
	validTypes := []string{"パブリック", "PC/CUG"}
	for _, t := range validTypes {
		if t == data {
			return nil
		}
	}

	return fmt.Errorf("invalid IXPeerType: %s (valid values: パブリック, PC/CUG)", data)
}

func CheckIXFields(ix, peerType, vlanID string) error {
	// IX名のチェック
	if err := CheckIncludeIXTemplate(ix); err != nil {
		return err
	}

	// IXPeerTypeのチェック
	if peerType == "" {
		return fmt.Errorf("IXPeerType is required for IXP connection")
	}
	if err := CheckIXPeerType(peerType); err != nil {
		return err
	}

	// PC/CUGの場合はVLAN-IDが必須
	if peerType == "PC/CUG" && vlanID == "" {
		return fmt.Errorf("IXVlanID is required when IXPeerType is PC/CUG")
	}

	return nil
}
