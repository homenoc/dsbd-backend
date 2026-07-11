package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func JoinPlan(ipID uint, input core.Plan) error {
	db := store.DB()

	return db.Model(&core.IP{Model: gorm.Model{ID: ipID}}).
		Association("Plan").
		Append(&input)
}

func DeletePlan(id uint) error {
	db := store.DB()

	return db.Delete(core.Plan{Model: gorm.Model{ID: id}}).Error
}

func UpdatePlan(input core.Plan) error {
	db := store.DB()

	return db.Model(&core.Plan{Model: gorm.Model{ID: input.ID}}).Updates(input).Error
}

func GetPlan(data *core.Plan) (core.Plan, error) {
	var plans core.Plan

	db := store.DB()

	err := db.First(&plans, data.ID).Error

	return plans, err
}
