package v0

import (
	"fmt"
	"log"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(g *core.Group) (*core.Group, error) {
	result := GetByOrg(g.Org)
	if result.Err != nil {
		return &core.Group{}, result.Err
	}
	if len(result.Group) != 0 {
		log.Println("error: this Org Name is already registered: " + g.Org)
		return &core.Group{}, fmt.Errorf("error: this org name is already registered")
	}

	err := store.DB().Create(&g).Error
	return g, err
}

func Delete(g *core.Group) error {
	return store.DB().Delete(g).Error
}

// UpdateAll updates the group with all non-zero fields of g (was Update(UpdateAll, ...)).
func UpdateAll(g core.Group) error {
	return store.DB().Model(&core.Group{Model: gorm.Model{ID: g.ID}}).Updates(g).Error
}

// GetByID loads a group with the full user/service/connection association graph.
func GetByID(id uint) group.ResultDatabase {
	var groups []core.Group
	err := store.DB().Preload("Users").
		Preload("Services").
		Preload("Tickets").
		Preload("Memos").
		Preload("Services.IP").
		Preload("Services.IP.Plan").
		Preload("Services.Connection").
		Preload("Services.Connection.BGPRouter").
		Preload("Services.Connection.BGPRouter.NOC").
		Preload("Services.Connection.TunnelEndPointRouterIP").
		Preload("Services.JPNICAdmin").
		Preload("Services.JPNICTech").
		First(&groups, id).Error
	return group.ResultDatabase{Group: groups, Err: err}
}

// GetByOrg returns groups matching an organization name.
func GetByOrg(org string) group.ResultDatabase {
	var groups []core.Group
	err := store.DB().Where("org = ?", org).Find(&groups).Error
	return group.ResultDatabase{Group: groups, Err: err}
}

func GetAll() group.ResultDatabase {
	var groups []core.Group
	err := store.DB().Preload("Users").
		Preload("Memos").
		Find(&groups).Error
	return group.ResultDatabase{Group: groups, Err: err}
}
