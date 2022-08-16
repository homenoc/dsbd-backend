package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	"gorm.io/gorm"
)

func userExtraction(inputUser, inputGroup, inputNOC []uint) []uint {
	var userArray []uint
	var groupArray []uint

	// if noc length isn't zero value
	if len(inputNOC) != 0 {

	}

	// if group length isn't zero value
	if len(inputGroup) != 0 {
		//I should implement check function
		for _, tmpGroup := range inputGroup {
			result := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: tmpGroup}})
			if result.Err != nil {
				for _, tmpResultGroup := range result.Group {
					for _, tmpUser := range tmpResultGroup.Users {
						userArray = append(userArray, tmpUser.ID)
					}
				}
			}
		}
		groupArray = removeDuplicate(groupArray)
	}

	// if user length isn't zero value
	if len(inputUser) != 0 {
		//I should implement check function
		for _, tmpUser := range inputUser {
			userArray = append(userArray, tmpUser)
		}
		userArray = removeDuplicate(userArray)
	}

	return userArray
}

func removeDuplicate(args []uint) []uint {
	results := make([]uint, 0, len(args))
	encountered := map[uint]bool{}
	for i := 0; i < len(args); i++ {
		if !encountered[args[i]] {
			encountered[args[i]] = true
			results = append(results, args[i])
		}
	}
	return results
}
