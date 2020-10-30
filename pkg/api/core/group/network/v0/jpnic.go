package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicTech"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbJpnicAdmin "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/jpnicAdmin/v0"
	dbJpnicTech "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/jpnicTech/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
)

type jpnic struct {
	admin   uint
	tech    []uint
	network group.Network
}

func jpnicProcess(input jpnic) error {
	log.Println(1)
	// input group id
	groupId := input.network.GroupID

	// 入力されたユーザのGroupIDを検索
	resultAdminUser := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: input.admin}})
	if resultAdminUser.Err != nil {
		return resultAdminUser.Err
	}
	// 認証ユーザのGroupIDと入力されたjpnicAdminの値が一致しているか確認
	if resultAdminUser.User[0].GID != groupId {
		return fmt.Errorf("This user's group id is not match. ")
	}

	// JPNIC Admin tableに保存
	_, err := dbJpnicAdmin.Create(&jpnicAdmin.JpnicAdmin{NetworkID: input.network.ID, UserID: input.admin, Lock: &[]bool{true}[0]})
	if err != nil {
		return err
	}

	//jpnic Tech Process
	for _, tmp := range input.tech {
		// 入力されたユーザのGroupIDを検索
		resultTechUser := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: tmp}})
		if resultTechUser.Err != nil {
			return resultTechUser.Err
		}

		// 認証ユーザのGroupIDと入力されたjpnicAdminの値が一致しているか確認
		if resultTechUser.User[0].GID != groupId {
			return fmt.Errorf("This user's group id is not match. ")
		}

		// JPNIC Tech tableに保存
		_, err := dbJpnicTech.Create(&jpnicTech.JpnicTech{NetworkID: input.network.ID, UserID: tmp, Lock: &[]bool{true}[0]})
		if err != nil {
			return err
		}
	}
	return nil
}
