package v0

import (
	"fmt"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	networkJPNICUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnic_user"
	networkUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/network_user"
	dbJPNICUser "github.com/homenoc/dsbd-backend/pkg/store/group/network/jpnic_user/v0"
	dbNetworkUser "github.com/homenoc/dsbd-backend/pkg/store/group/network/network_user/v0"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/store/group/network/v0"
	"github.com/jinzhu/gorm"
)

func checkGroupID(groupID, NetworkID, JPNICUserID uint) error {
	resultNetwork := dbNetwork.Get(network.ID, &network.Network{Model: gorm.Model{ID: NetworkID}})
	if resultNetwork.Err != nil {
		return resultNetwork.Err
	}
	if resultNetwork.Network[0].GroupID != groupID {
		return fmt.Errorf("not match Network GroupID")
	}
	resultJPNICUser := dbJPNICUser.Get(networkJPNICUser.ID, &networkJPNICUser.JPNICUser{Model: gorm.Model{ID: JPNICUserID}})
	if resultJPNICUser.Err != nil {
		return resultJPNICUser.Err
	}
	if resultNetwork.Network[0].GroupID != groupID {
		return fmt.Errorf("not match JPNICUser GroupID")
	}
	return nil
}

func checkDuplicate(input networkUser.NetworkUser) error {
	result := dbNetworkUser.Get(networkUser.Info, &input)
	if result.Err != nil {
		return result.Err
	}
	if len(result.NetworkUser) > 0 {
		return fmt.Errorf("duplicate")
	}
	return nil
}
