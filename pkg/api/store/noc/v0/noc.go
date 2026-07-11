package v0

import (
	core "github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(n *core.NOC) (*core.NOC, error) {
	err := store.DB().Create(&n).Error
	return n, err
}

func Delete(n *core.NOC) error {
	return store.DB().Delete(n).Error
}

// UpdateAll updates the editable NOC fields (was Update(UpdateAll, ...)).
func UpdateAll(data core.NOC) error {
	return store.DB().Model(&core.NOC{Model: gorm.Model{ID: data.ID}}).Updates(core.NOC{
		Name:      data.Name,
		Location:  data.Location,
		Bandwidth: data.Bandwidth,
		Enable:    data.Enable,
		Comment:   data.Comment,
	}).Error
}

// GetByID looks up one NOC by primary key.
func GetByID(id uint) noc.ResultDatabase {
	var nocs []core.NOC
	err := store.DB().First(&nocs, id).Error
	return noc.ResultDatabase{NOC: nocs, Err: err}
}

func GetAll() noc.ResultDatabase {
	var nocs []core.NOC
	err := store.DB().Preload("BGPRouter").
		Preload("TunnelEndPointRouter").
		Preload("TunnelEndPointRouter.TunnelEndPointRouterIP").
		Find(&nocs).Error
	return noc.ResultDatabase{NOC: nocs, Err: err}
}
