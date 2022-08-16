package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(g *core.Group) (*core.Group, error) {
	result := Get(group.Org, &core.Group{Org: g.Org})
	if result.Err != nil {
		return &core.Group{}, result.Err
	}
	if len(result.Group) != 0 {
		log.Println("error: this Org Name is already registered: " + g.Org)
		return &core.Group{}, fmt.Errorf("error: this org name is already registered")
	}

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return g, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&g).Error
	return g, err
}

func Delete(group *core.Group) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	return db.Delete(group).Error
}

func Update(base int, g core.Group) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = nil

	if group.UpdateOrg == base {
		err = db.Model(&core.Group{Model: gorm.Model{ID: g.ID}}).Updates(core.Group{Org: g.Org}).Error
	} else if group.UpdateMembership == base {
		err = db.Model(&core.Group{Model: gorm.Model{ID: g.ID}}).Updates(core.Group{
			StripeCustomerID:     g.StripeCustomerID,
			StripeSubscriptionID: g.StripeSubscriptionID,
			MemberExpired:        g.MemberExpired,
		}).Error
	} else if group.UpdateAll == base {
		err = db.Model(&core.Group{Model: gorm.Model{ID: g.ID}}).Updates(g).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return err
}

func Get(base int, data *core.Group) group.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var groupStruct []core.Group

	if base == group.ID { //ID
		err = db.Preload("Users").
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
			First(&groupStruct, data.ID).Error
	} else if base == group.Org { //Org
		err = db.Where("org = ?", data.Org).Find(&groupStruct).Error
	} else {
		log.Println("base select error")
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return group.ResultDatabase{Group: groupStruct, Err: err}
}

func GetAll() group.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var groups []core.Group
	err = db.Preload("Users").
		Preload("Memos").
		Find(&groups).Error
	return group.ResultDatabase{Group: groups, Err: err}
}
