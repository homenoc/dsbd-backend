package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/ip"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(network *core.IP) (*core.IP, error) {
	err := store.DB().Create(&network).Error
	return network, err
}

func Delete(network *core.IP) error {
	return store.DB().Delete(network).Error
}

// UpdateAll updates the IP with all non-zero fields of u (was Update(UpdateAll, ...)).
func UpdateAll(u core.IP) error {
	return store.DB().Model(&core.IP{Model: gorm.Model{ID: u.ID}}).Updates(u).Error
}

// GetByID looks up one IP by primary key.
func GetByID(id uint) ip.ResultDatabase {
	var ips []core.IP
	err := store.DB().First(&ips, id).Error
	return ip.ResultDatabase{IP: ips, Err: err}
}

func GetAll() ip.ResultDatabase {
	var ips []core.IP
	err := store.DB().Find(&ips).Error
	return ip.ResultDatabase{IP: ips, Err: err}
}
