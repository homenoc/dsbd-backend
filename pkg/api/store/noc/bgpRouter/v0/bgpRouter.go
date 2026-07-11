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

// UpdateAll updates the editable BGP router fields (was Update(UpdateAll, ...)).
func UpdateAll(data core.BGPRouter) error {
	return store.DB().Model(&core.BGPRouter{Model: gorm.Model{ID: data.ID}}).Updates(core.BGPRouter{
		HostName: data.HostName,
		Address:  data.Address,
		Enable:   data.Enable,
	}).Error
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
