package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func JoinJPNICTech(input core.JPNICTech) error {
	db := store.DB()

	return db.Create(&input).Error
}

func DeleteJPNICTech(id uint) error {
	db := store.DB()

	return db.Delete(core.JPNICTech{Model: gorm.Model{ID: id}}).Error
}

func UpdateJPNICTech(input core.JPNICTech) error {
	db := store.DB()

	return db.Model(&core.JPNICTech{Model: gorm.Model{ID: input.ID}}).Updates(input).Error
}
