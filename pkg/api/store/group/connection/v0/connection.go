package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(c *core.Connection) (*core.Connection, error) {
	err := store.DB().Create(&c).Error
	return c, err
}

func Delete(c *core.Connection) error {
	return store.DB().Delete(c).Error
}

// UpdateAll updates every connection field (Select("*")), preserving the
// previous Update(UpdateAll, ...) behavior including zero-valued fields.
func UpdateAll(c core.Connection) error {
	return store.DB().Model(&core.Connection{Model: gorm.Model{ID: c.ID}}).Select("*").Updates(c).Error
}

// GetByID loads one connection with router/service/group associations.
func GetByID(id uint) connection.ResultDatabase {
	var connections []core.Connection
	err := store.DB().Preload("BGPRouter").
		Preload("BGPRouter.NOC").
		Preload("TunnelEndPointRouterIP").
		Preload("Service").
		Preload("Service.Group").
		First(&connections, id).Error
	return connection.ResultDatabase{Connection: connections, Err: err}
}

// GetByServiceID returns the connections belonging to a service.
func GetByServiceID(serviceID uint) connection.ResultDatabase {
	var connections []core.Connection
	err := store.DB().Where("service_id = ?", serviceID).Find(&connections).Error
	return connection.ResultDatabase{Connection: connections, Err: err}
}

func GetAll() connection.ResultDatabase {
	var connections []core.Connection
	err := store.DB().Preload("BGPRouter").
		Preload("BGPRouter.NOC").
		Preload("TunnelEndPointRouterIP").
		Preload("Service").
		Preload("Service.Group").
		Find(&connections).Error
	return connection.ResultDatabase{Connection: connections, Err: err}
}
