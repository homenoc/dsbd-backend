package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func JoinPlan(ipID uint, input core.Plan) error {
	return store.DB().Model(&core.IP{Model: gorm.Model{ID: ipID}}).
		Association("Plan").
		Append(&input)
}

func DeletePlan(id uint) error {
	return store.DB().Delete(core.Plan{Model: gorm.Model{ID: id}}).Error
}

// UpdatePlan writes the editable plan columns from a full object. All are
// value-typed and always written (clearing/zeroing persists); ip_id stays
// untouchable.
func UpdatePlan(input core.Plan) error {
	cols := []string{"name", "after", "half_year", "one_year"}
	return store.DB().Model(&core.Plan{Model: gorm.Model{ID: input.ID}}).Select(cols).Updates(input).Error
}

func GetPlan(data *core.Plan) (core.Plan, error) {
	var plans core.Plan

	err := store.DB().First(&plans, data.ID).Error

	return plans, err
}
