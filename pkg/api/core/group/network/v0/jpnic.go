package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"log"
)

type jpnicHandler struct {
	admin      uint
	tech       []uint
	groupID    uint
	jpnicAdmin *network.JPNICAdmin
	jpnicTech  *[]network.JPNICTech
}

func (jpnic *jpnicHandler) jpnicProcess() error {
	log.Println(jpnic)
	// 入力されたユーザのGroupIDを検索
	resultGroupUser := dbUser.Get(user.GID, &user.User{GroupID: jpnic.groupID})
	if resultGroupUser.Err != nil {
		return resultGroupUser.Err
	}

	// 管理者連絡窓口
	for _, tmpUser := range resultGroupUser.User {
		if tmpUser.ID == jpnic.admin {
			jpnic.jpnicAdmin = &network.JPNICAdmin{UserID: jpnic.admin, Lock: &[]bool{true}[0]}
			break
		}
	}

	// groupIDに対してJPNICAdminが見つからなかった場合
	if jpnic.jpnicAdmin == nil {
		return fmt.Errorf("This user have no authorization. ")
	}

	// 技術連絡担当者
	var jpnicTech []network.JPNICTech
	for _, tmpTechUserID := range jpnic.tech {
		// ユーザの権限確認
		for _, tmpUser := range resultGroupUser.User {
			if tmpUser.ID == tmpTechUserID {
				jpnicTech = append(jpnicTech, network.JPNICTech{UserID: jpnic.admin, Lock: &[]bool{true}[0]})
				break
			}
		}
	}

	if jpnicTech == nil {
		return fmt.Errorf("This user have no authorization. ")
	}

	jpnic.jpnicTech = &jpnicTech

	return nil
}

//func jpnicProcess(input jpnic) error {
//	log.Println(1)
//	// input group id
//	groupId := input.network.GroupID
//
//	// 入力されたユーザのGroupIDを検索
//	resultAdminUser := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: input.admin}})
//	if resultAdminUser.Err != nil {
//		return resultAdminUser.Err
//	}
//	// 認証ユーザのGroupIDと入力されたjpnicAdminの値が一致しているか確認
//	if resultAdminUser.User[0].GroupID != groupId {
//		return fmt.Errorf("This user's group id is not match. ")
//	}
//
//	// JPNIC Admin tableに保存
//	_, err := dbJpnicAdmin.Create(&jpnicAdmin.JpnicAdmin{NetworkID: input.network.ID, UserID: input.admin, Lock: &[]bool{true}[0]})
//	if err != nil {
//		return err
//	}
//
//	//jpnic Tech Process
//	for _, tmp := range input.tech {
//		// 入力されたユーザのGroupIDを検索
//		resultTechUser := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: tmp}})
//		if resultTechUser.Err != nil {
//			return resultTechUser.Err
//		}
//
//		// 認証ユーザのGroupIDと入力されたjpnicAdminの値が一致しているか確認
//		if resultTechUser.User[0].GroupID != groupId {
//			return fmt.Errorf("This user's group id is not match. ")
//		}
//
//		// JPNIC Tech tableに保存
//		_, err := dbJpnicTech.Create(&jpnicTech.JpnicTech{NetworkID: input.network.ID, UserID: tmp, Lock: &[]bool{true}[0]})
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
