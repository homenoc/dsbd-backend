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

func UpdateIP(input core.IP) error {
	return store.DB().Model(&core.IP{Model: gorm.Model{ID: input.ID}}).Updates(input).Error
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
