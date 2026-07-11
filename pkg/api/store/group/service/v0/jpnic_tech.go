package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func JoinJPNICTech(input core.JPNICTech) error {
	return store.DB().Create(&input).Error
}

func DeleteJPNICTech(id uint) error {
	return store.DB().Delete(core.JPNICTech{Model: gorm.Model{ID: id}}).Error
}

func UpdateJPNICTech(input core.JPNICTech) error {
	return store.DB().Model(&core.JPNICTech{Model: gorm.Model{ID: input.ID}}).Updates(input).Error
}

func GetJPNICTech(id uint) (core.JPNICTech, error) {
	var tech core.JPNICTech
	err := store.DB().First(&tech, id).Error
	return tech, err
}
