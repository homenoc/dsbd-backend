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

// UpdateJPNICTech writes the editable JPNIC tech-contact columns from a full
// object. Value-typed columns are always written (clearing persists);
// service_id and the hidden/is_group structural flags stay untouchable.
func UpdateJPNICTech(input core.JPNICTech) error {
	cols := []string{"hidden", "is_group", "v4_jpnic_handle", "v6_jpnic_handle", "name", "name_en", "mail",
		"org", "org_en", "post_code", "address", "address_en",
		"dept", "dept_en", "title", "title_en", "tel", "fax", "country"}
	return store.DB().Model(&core.JPNICTech{Model: gorm.Model{ID: input.ID}}).
		Select(cols).Updates(input).Error
}

func GetJPNICTech(id uint) (core.JPNICTech, error) {
	var tech core.JPNICTech
	err := store.DB().First(&tech, id).Error
	return tech, err
}
