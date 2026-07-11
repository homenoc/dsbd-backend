package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(s *core.Service) (*core.Service, error) {
	err := store.DB().Create(&s).Error
	return s, err
}

func Delete(s *core.Service) error {
	return store.DB().Delete(s).Error
}

// UpdateData updates the editable applicant/org fields (was Update(UpdateData, ...)).
func UpdateData(c core.Service) error {
	return store.DB().Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Updates(core.Service{
		Org:       c.Org,
		OrgEn:     c.OrgEn,
		PostCode:  c.PostCode,
		Address:   c.Address,
		AddressEn: c.AddressEn,
		ASN:       c.ASN,
	}).Error
}

// UpdateAll updates the service with all non-zero fields of c (was Update(UpdateAll, ...)).
func UpdateAll(c core.Service) error {
	return store.DB().Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Updates(c).Error
}

// GetByID loads one service with its full association graph.
func GetByID(id uint) service.ResultDatabase {
	var services []core.Service
	err := store.DB().Preload("IP").
		Preload("IP.Plan").
		Preload("Connection").
		Preload("Connection.BGPRouter").
		Preload("Connection.TunnelEndPointRouterIP").
		Preload("JPNICAdmin").
		Preload("JPNICTech").
		Preload("Group").
		First(&services, id).Error
	return service.ResultDatabase{Err: err, Service: services}
}

// GetByGroupID returns a group's services (used for numbering new services).
func GetByGroupID(groupID uint) service.ResultDatabase {
	var services []core.Service
	err := store.DB().Where("group_id = ?", groupID).Find(&services).Error
	return service.ResultDatabase{Err: err, Service: services}
}

// GetAddAllowByGroupID returns a group's services that currently allow additions.
func GetAddAllowByGroupID(groupID uint) service.ResultDatabase {
	var services []core.Service
	err := store.DB().Where("group_id = ? AND add_allow = ?", groupID, true).Find(&services).Error
	return service.ResultDatabase{Err: err, Service: services}
}

// GetByASN returns passed+enabled services for an ASN with open connections
// (used by Slack ASN lookup).
func GetByASN(asn *uint) service.ResultDatabase {
	var services []core.Service
	err := store.DB().Where("asn = ? AND pass = ? AND enable = ?", asn, true, true).
		Preload("IP").
		Preload("Connection", "open = ?", true).
		Preload("Group").
		Find(&services).Error
	return service.ResultDatabase{Err: err, Service: services}
}

func GetAll() service.ResultDatabase {
	var services []core.Service
	err := store.DB().Preload("IP").
		Preload("Connection").
		Preload("Connection.BGPRouter").
		Preload("Connection.TunnelEndPointRouterIP").
		Preload("JPNICAdmin").
		Preload("JPNICTech").
		Find(&services).Error
	return service.ResultDatabase{Err: err, Service: services}
}
