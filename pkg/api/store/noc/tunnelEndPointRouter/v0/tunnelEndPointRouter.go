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

// Update writes the admin-editable columns of the tunnel endpoint router row.
// Value-typed columns are always written (so clearing persists); pointer
// columns only when provided.
func Update(data core.TunnelEndPointRouter) error {
	cols := []string{"host_name", "capacity", "comment"}
	if data.NOCID != nil {
		cols = append(cols, "noc_id")
	}
	if data.Enable != nil {
		cols = append(cols, "enable")
	}
	return store.DB().Model(&core.TunnelEndPointRouter{Model: gorm.Model{ID: data.ID}}).
		Select(cols).Updates(data).Error
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
