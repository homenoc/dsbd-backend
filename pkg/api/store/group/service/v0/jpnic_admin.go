package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func JoinJPNICByAdmin(serviceID uint, input core.JPNICAdmin) error {
	return store.DB().Model(&core.Service{Model: gorm.Model{ID: serviceID}}).
		Association("JPNICAdmin").
		Append(&input)
}

func DeleteJPNICByAdmin(id uint) error {
	return store.DB().Delete(core.JPNICAdmin{Model: gorm.Model{ID: id}}).Error
}

func UpdateJPNICByAdmin(input core.JPNICAdmin) error {
	return store.DB().Model(&core.JPNICAdmin{Model: gorm.Model{ID: input.ID}}).Updates(input).Error
}
