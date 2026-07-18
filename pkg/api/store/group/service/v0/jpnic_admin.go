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

// UpdateJPNICByAdmin writes the editable JPNIC admin-contact columns from a
// full object. Value-typed columns are always written (clearing persists,
// and hidden/is_group are settable both ways); service_id stays untouchable.
func UpdateJPNICByAdmin(input core.JPNICAdmin) error {
	cols := []string{"hidden", "is_group", "v4_jpnic_handle", "v6_jpnic_handle", "name", "name_en", "mail",
		"org", "org_en", "post_code", "address", "address_en",
		"dept", "dept_en", "title", "title_en", "tel", "fax", "country"}
	return store.DB().Model(&core.JPNICAdmin{Model: gorm.Model{ID: input.ID}}).
		Select(cols).Updates(input).Error
}

func GetJPNICAdmin(id uint) (core.JPNICAdmin, error) {
	var admin core.JPNICAdmin
	err := store.DB().First(&admin, id).Error
	return admin, err
}
