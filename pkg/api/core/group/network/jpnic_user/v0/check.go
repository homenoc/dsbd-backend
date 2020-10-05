package v0

import (
	"fmt"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnic_user"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/store/group/network/v0"
	"github.com/jinzhu/gorm"
)

func check(input jpnic_user.JPNICUser) error {
	// check
	if input.NameJa == "" {
		return fmt.Errorf("no data: name japanese")
	}
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	if input.OrgJa == "" {
		return fmt.Errorf("no data: position japanese")
	}
	if input.Org == "" {
		return fmt.Errorf("no data: position")
	}
	if input.PostCode == "" {
		return fmt.Errorf("no data: postcode")
	}
	if input.AddressJa == "" {
		return fmt.Errorf("no data: address japanese")
	}
	if input.Address == "" {
		return fmt.Errorf("no data: address")
	}
	if input.Mail == "" {
		return fmt.Errorf("no data: mail")
	}
	if input.Tel == "" {
		return fmt.Errorf("no data: tel")
	}
	return nil
}

func checkNetworkID(input jpnic_user.JPNICUser, groupID uint) error {
	if input.OperationID == 0 && input.TechID == 0 {
		return nil
	}
	if input.OperationID != 0 {
		result := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: input.OperationID}})
		if len(result.Network) == 0 {
			return fmt.Errorf("mistake operation ID")
		}
		if result.Network[0].GroupID == groupID {
			return nil
		} else {
			return fmt.Errorf("mismatch: Group ID| OperationID")
		}
	}

	if input.TechID != 0 {
		result := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: input.TechID}})
		if len(result.Network) == 0 {
			return fmt.Errorf("mistake tech ID")
		}
		if result.Network[0].GroupID == groupID {
			return nil
		} else {
			return fmt.Errorf("mismatch: Group ID| TechID")
		}
	}
	return nil
}
