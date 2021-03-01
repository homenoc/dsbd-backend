package v0

import (
	"fmt"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
)

func check(groupID uint, restrict bool, input connection.Connection) error {
	// check

	resultNOC := dbNOC.GetAll()
	if resultNOC.Err != nil {
		return resultNOC.Err
	}

	exists := false

	for _, tmp := range resultNOC.NOC {
		if input.NOC == tmp.Name {
			exists = true
			break
		}
		if input.ConnectionType == "any" {
			exists = true
			break
		}
	}

	if !exists {
		return fmt.Errorf("no data: connection Type")
	}

	if input.ConnectionType == "EIP" || input.ConnectionType == "IPT" {
		if input.NTT == "" {
			return fmt.Errorf("no data: ntt")
		}

		if input.TermIP == "" {
			return fmt.Errorf("no data: Term IP")
		}
	}

	if input.ConnectionType == "CC0" || input.ConnectionType == "ETC" {
		if input.ConnectionComment == "" {
			return fmt.Errorf("no data: connection comment")
		}
	}

	if input.UserID == 0 {
		return fmt.Errorf("no data: userID")
	}

	if restrict {
		resultUser := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: input.UserID}})
		if resultUser.Err != nil {
			return resultUser.Err
		}

		log.Println(groupID)
		log.Println(resultUser.User[0].GroupID)

		if groupID != resultUser.User[0].GroupID {
			return fmt.Errorf("error: not match groupID")
		}
	}

	return nil
}
