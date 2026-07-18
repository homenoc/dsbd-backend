package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func CreateIP(input core.IP) error {
	return store.DB().Create(&input).Error
}

func DeleteIP(id uint) error {
	return store.DB().Select("Plan").Delete(&core.IP{Model: gorm.Model{ID: id}}).Error
}

// UpdateIP writes the admin-editable IP columns from a full object. Value-typed
// columns are always written (clearing persists); pointer/time columns only
// when provided. service_id/version stay untouchable.
func UpdateIP(input core.IP) error {
	cols := []string{"name", "ip", "use_case"}
	if !input.StartDate.IsZero() {
		cols = append(cols, "start_date")
	}
	if input.EndDate != nil {
		cols = append(cols, "end_date")
	}
	if input.Open != nil {
		cols = append(cols, "open")
	}
	return store.DB().Model(&core.IP{Model: gorm.Model{ID: input.ID}}).Select(cols).Updates(input).Error
}

func GetIP(id uint) (core.IP, error) {
	var ip core.IP
	err := store.DB().First(&ip, id).Error
	return ip, err
}

func GetAllIP() ([]core.IP, error) {
	var ips []core.IP
	err := store.DB().Find(&ips).Error
	return ips, err
}
