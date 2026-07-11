package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouter"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(r *core.TunnelEndPointRouter) (*core.TunnelEndPointRouter, error) {
	err := store.DB().Create(&r).Error
	return r, err
}

func Delete(r *core.TunnelEndPointRouter) error {
	return store.DB().Delete(r).Error
}

// UpdateAll updates the editable tunnel endpoint router fields.
func UpdateAll(data core.TunnelEndPointRouter) error {
	return store.DB().Model(&core.TunnelEndPointRouter{Model: gorm.Model{ID: data.ID}}).
		Updates(core.TunnelEndPointRouter{
			NOCID:    data.NOCID,
			HostName: data.HostName,
			Capacity: data.Capacity,
			Comment:  data.Comment,
			Enable:   data.Enable,
		}).Error
}

// GetByID looks up one tunnel endpoint router by primary key.
func GetByID(id uint) tunnelEndPointRouter.ResultDatabase {
	var routers []core.TunnelEndPointRouter
	err := store.DB().First(&routers, id).Error
	return tunnelEndPointRouter.ResultDatabase{TunnelEndPointRouter: routers, Err: err}
}

func GetAll() tunnelEndPointRouter.ResultDatabase {
	var routers []core.TunnelEndPointRouter
	err := store.DB().Find(&routers).Error
	return tunnelEndPointRouter.ResultDatabase{TunnelEndPointRouter: routers, Err: err}
}
