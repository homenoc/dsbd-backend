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

// Update writes the admin-editable columns of the connection row from a full
// object (the admin whole-object PUT). This replaces Select("*"), which nulled
// every omitted column — including service_id — on a sparse payload. Value-typed
// columns are always written (clearing to "" persists); pointer columns only
// when provided. Consequence: bgp_router_id / tunnel_end_point_router_ip_id can
// no longer be nulled by omission — the FE's "なし" convention writes 0 instead.
func Update(c core.Connection) error {
	cols := []string{"connection_type", "connection_comment", "ix", "ix_peer_type",
		"ix_vlan_id", "ipv4_route", "ipv6_route", "ntt", "preferred_ap", "term_ip",
		"rfc8950", "address", "link_v4_our", "link_v4_your", "link_v6_our",
		"link_v6_your", "comment"}
	if c.BGPRouterID != nil {
		cols = append(cols, "bgp_router_id")
	}
	if c.TunnelEndPointRouterIPID != nil {
		cols = append(cols, "tunnel_end_point_router_ip_id")
	}
	if c.Monitor != nil {
		cols = append(cols, "monitor")
	}
	if c.Open != nil {
		cols = append(cols, "open")
	}
	if c.Enable != nil {
		cols = append(cols, "enable")
	}
	return store.DB().Model(&core.Connection{Model: gorm.Model{ID: c.ID}}).Select(cols).Updates(c).Error
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
