package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func JoinIP(input core.IP) error {
	db := store.DB()

	return db.Create(&input).Error
}

func DeleteIP(id uint) error {
	db := store.DB()

	return db.Select("Plan").Delete(&core.IP{Model: gorm.Model{ID: id}}).Error
}

func UpdateIP(input core.IP) error {
	db := store.DB()

	return db.Model(&core.IP{Model: gorm.Model{ID: input.ID}}).Updates(input).Error
}
