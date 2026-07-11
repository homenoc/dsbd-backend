package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/jpnicAdmin"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(network *core.JPNICAdmin) (*core.JPNICAdmin, error) {
	err := store.DB().Create(&network).Error
	return network, err
}

func Delete(network *core.JPNICAdmin) error {
	return store.DB().Delete(network).Error
}

// UpdateAll updates the editable JPNIC admin fields (was Update(UpdateAll, ...)).
func UpdateAll(u core.JPNICAdmin) error {
	return store.DB().Model(&core.JPNICAdmin{Model: gorm.Model{ID: u.ID}}).Updates(core.JPNICAdmin{
		Org:       u.Org,
		OrgEn:     u.OrgEn,
		PostCode:  u.PostCode,
		Address:   u.Address,
		AddressEn: u.AddressEn,
		Dept:      u.Dept,
		DeptEn:    u.DeptEn,
		Tel:       u.Tel,
		Fax:       u.Fax,
		Country:   u.Country,
	}).Error
}

// GetByID looks up one JPNIC admin record by primary key.
func GetByID(id uint) jpnicAdmin.ResultDatabase {
	var admins []core.JPNICAdmin
	err := store.DB().First(&admins, id).Error
	return jpnicAdmin.ResultDatabase{Admins: admins, Err: err}
}

func GetAll() jpnicAdmin.ResultDatabase {
	var admins []core.JPNICAdmin
	err := store.DB().Find(&admins).Error
	return jpnicAdmin.ResultDatabase{Admins: admins, Err: err}
}
