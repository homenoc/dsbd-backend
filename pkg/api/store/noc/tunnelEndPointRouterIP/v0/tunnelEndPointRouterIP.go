package v0

import (
	core "github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouterIP"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(r *core.TunnelEndPointRouterIP) (*core.TunnelEndPointRouterIP, error) {
	err := store.DB().Create(&r).Error
	return r, err
}

func Delete(r *core.TunnelEndPointRouterIP) error {
	return store.DB().Delete(r).Error
}

// Update writes the admin-editable columns of the tunnel endpoint router IP
// row. Value-typed columns are always written (so clearing persists); pointer
// columns only when provided.
func Update(data core.TunnelEndPointRouterIP) error {
	cols := []string{"ip", "comment"}
	if data.Enable != nil {
		cols = append(cols, "enable")
	}
	return store.DB().Model(&core.TunnelEndPointRouterIP{Model: gorm.Model{ID: data.ID}}).
		Select(cols).Updates(data).Error
}

// GetByID looks up one tunnel endpoint router IP by primary key (with its router).
func GetByID(id uint) tunnelEndPointRouterIP.ResultDatabase {
	var routers []core.TunnelEndPointRouterIP
	err := store.DB().Preload("TunnelEndPointRouter").First(&routers, id).Error
	return tunnelEndPointRouterIP.ResultDatabase{TunnelEndPointRouterIP: routers, Err: err}
}

func GetAll() tunnelEndPointRouterIP.ResultDatabase {
	var routers []core.TunnelEndPointRouterIP
	err := store.DB().Preload("TunnelEndPointRouter").Find(&routers).Error
	return tunnelEndPointRouterIP.ResultDatabase{TunnelEndPointRouterIP: routers, Err: err}
}
