package core

import "fmt"

func GetMembershipTypeID(id uint) (ConstantMembership, error) {
	for _, memberType := range MemberTypes {
		if memberType.ID == id {
			return memberType, nil
		}
	}
	return ConstantMembership{}, fmt.Errorf("error: getting membership")
}
