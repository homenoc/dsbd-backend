package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/jpnicTech"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(jpnic *core.JPNICTech) (*core.JPNICTech, error) {
	err := store.DB().Create(&jpnic).Error
	return jpnic, err
}

func Delete(jpnic *core.JPNICTech) error {
	return store.DB().Delete(jpnic).Error
}

// UpdateAll updates the JPNIC tech record with all non-zero fields of jpnic.
func UpdateAll(jpnic core.JPNICTech) error {
	return store.DB().Model(&core.JPNICTech{Model: gorm.Model{ID: jpnic.ID}}).Updates(jpnic).Error
}

// GetByID looks up one JPNIC tech record by primary key.
func GetByID(id uint) jpnicTech.ResultDatabase {
	var techs []core.JPNICTech
	err := store.DB().First(&techs, id).Error
	return jpnicTech.ResultDatabase{Tech: techs, Err: err}
}

func GetAll() jpnicTech.ResultDatabase {
	var techs []core.JPNICTech
	err := store.DB().Find(&techs).Error
	return jpnicTech.ResultDatabase{Tech: techs, Err: err}
}
