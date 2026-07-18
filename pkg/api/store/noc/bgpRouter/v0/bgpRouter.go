package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	router "github.com/homenoc/dsbd-backend/pkg/api/core/noc/bgpRouter"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(r *core.BGPRouter) (*core.BGPRouter, error) {
	err := store.DB().Create(&r).Error
	return r, err
}

func Delete(r *core.BGPRouter) error {
	return store.DB().Delete(r).Error
}

// Update writes the admin-editable columns of the BGP router row. Value-typed
// columns are always written (so clearing to "" persists); pointer columns only
// when provided. Note: noc_id and comment were not updatable before either —
// kept out of the whitelist for parity.
func Update(data core.BGPRouter) error {
	cols := []string{"host_name", "address"}
	if data.Enable != nil {
		cols = append(cols, "enable")
	}
	return store.DB().Model(&core.BGPRouter{Model: gorm.Model{ID: data.ID}}).Select(cols).Updates(data).Error
}

// GetByID looks up one BGP router by primary key.
func GetByID(id uint) router.ResultDatabase {
	var routers []core.BGPRouter
	err := store.DB().First(&routers, id).Error
	return router.ResultDatabase{BGPRouter: routers, Err: err}
}

func GetAll() router.ResultDatabase {
	var routers []core.BGPRouter
	err := store.DB().Find(&routers).Error
	return router.ResultDatabase{BGPRouter: routers, Err: err}
}
