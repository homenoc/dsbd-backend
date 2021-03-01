package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"log"
)

type adminTechHandler struct {
	admin       uint
	tech        []uint
	groupID     uint
	resultAdmin *network.Admin
	resultTech  *[]network.Tech
}

func (adminTech *adminTechHandler) AdminTechProcess() error {
	log.Println(adminTech)
	// 入力されたユーザのGroupIDを検索
	resultGroupUser := dbUser.Get(user.GID, &user.User{GroupID: adminTech.groupID})
	if resultGroupUser.Err != nil {
		return resultGroupUser.Err
	}

	// 管理者連絡窓口
	for _, tmpUser := range resultGroupUser.User {
		if tmpUser.ID == adminTech.admin {
			adminTech.resultAdmin = &network.Admin{UserID: adminTech.admin, Lock: &[]bool{true}[0]}
			break
		}
	}

	// groupIDに対してJPNICAdminが見つからなかった場合
	if adminTech.resultAdmin == nil {
		return fmt.Errorf("This user have no authorization. ")
	}

	// 技術連絡担当者
	var tech []network.Tech
	for _, tmpTechUserID := range adminTech.tech {
		// ユーザの権限確認
		for _, tmpUser := range resultGroupUser.User {
			if tmpUser.ID == tmpTechUserID {
				tech = append(tech, network.Tech{UserID: adminTech.admin, Lock: &[]bool{true}[0]})
				break
			}
		}
	}

	if tech == nil {
		return fmt.Errorf("This user have no authorization. ")
	}

	adminTech.resultTech = &tech

	return nil
}
