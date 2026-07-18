package v0

import (
	"fmt"
	"log"
	"time"

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

// Update writes the admin-editable columns of the group row from a full
// object (the admin whole-object PUT). Value-typed columns are always written
// (so clearing to "" persists); pointer/time columns only when provided, so an
// omitted field never nulls the row. Sparse internal writes must use the
// intent functions below instead — this function would clear the value-typed
// columns a sparse struct leaves empty.
func Update(g core.Group) error {
	cols := []string{"question", "org", "org_en", "post_code", "address", "address_en",
		"tel", "country", "contract", "comment"}
	if g.MemberType != 0 {
		cols = append(cols, "member_type")
	}
	if g.Agree != nil {
		cols = append(cols, "agree")
	}
	if g.CouponID != nil {
		cols = append(cols, "coupon_id")
	}
	if g.MemberExpired != nil {
		cols = append(cols, "member_expired")
	}
	if g.Pass != nil {
		cols = append(cols, "pass")
	}
	if g.ExpiredStatus != nil {
		cols = append(cols, "expired_status")
	}
	if g.AddAllow != nil {
		cols = append(cols, "add_allow")
	}
	if g.StripeCustomerID != nil {
		cols = append(cols, "stripe_customer_id")
	}
	if g.StripeSubscriptionID != nil {
		cols = append(cols, "stripe_subscription_id")
	}
	return store.DB().Model(&core.Group{Model: gorm.Model{ID: g.ID}}).Select(cols).Updates(g).Error
}

// UpdateAddAllow flips whether the group may file new service applications.
func UpdateAddAllow(id uint, allow bool) error {
	return store.DB().Model(&core.Group{Model: gorm.Model{ID: id}}).
		Update("add_allow", allow).Error
}

// UpdateStripeCustomerID records the group's Stripe customer.
func UpdateStripeCustomerID(id uint, customerID string) error {
	return store.DB().Model(&core.Group{Model: gorm.Model{ID: id}}).
		Update("stripe_customer_id", customerID).Error
}

// UpdateStripeSubscription records a subscription and the resulting membership expiry.
func UpdateStripeSubscription(id uint, subscriptionID string, memberExpired time.Time) error {
	return store.DB().Model(&core.Group{Model: gorm.Model{ID: id}}).
		Updates(map[string]any{"stripe_subscription_id": subscriptionID, "member_expired": memberExpired}).Error
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
