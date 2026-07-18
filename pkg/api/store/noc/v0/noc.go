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

// Update writes the admin-editable columns of the NOC row. Value-typed columns
// are always written (so clearing to "" persists); pointer columns only when
// provided, so an omitted field never nulls the row. Callers pass a full object
// (the handler merges the request onto the current record).
func Update(data core.NOC) error {
	cols := []string{"name", "location", "bandwidth", "comment"}
	if data.Enable != nil {
		cols = append(cols, "enable")
	}
	return store.DB().Model(&core.NOC{Model: gorm.Model{ID: data.ID}}).Select(cols).Updates(data).Error
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
