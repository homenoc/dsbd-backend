package v0

import (
	"fmt"
	core "github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouterIP"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(router *core.TunnelEndPointRouterIP) (*core.TunnelEndPointRouterIP, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return router, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&router).Error
	return router, err
}

func Delete(router *core.TunnelEndPointRouterIP) error {
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

	return db.Delete(router).Error
}

func Update(base int, data core.TunnelEndPointRouterIP) error {
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

	if tunnelEndPointRouterIP.UpdateAll == base {
		err = db.Model(&core.TunnelEndPointRouterIP{Model: gorm.Model{ID: data.ID}}).Updates(core.TunnelEndPointRouterIP{
			IP:      data.IP,
			Comment: data.Comment,
			Enable:  data.Enable,
		}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return err
}

func Get(base int, data *core.TunnelEndPointRouterIP) tunnelEndPointRouterIP.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return tunnelEndPointRouterIP.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return tunnelEndPointRouterIP.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var routerStruct []core.TunnelEndPointRouterIP

	if base == tunnelEndPointRouterIP.ID { //ID
		err = db.Preload("TunnelEndPointRouter").
			First(&routerStruct, data.ID).Error
	} else if base == tunnelEndPointRouterIP.Enable { //GroupID
		err = db.Where("enable = ?", data.Enable).Find(&routerStruct).Error
	} else {
		log.Println("base select error")
		return tunnelEndPointRouterIP.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return tunnelEndPointRouterIP.ResultDatabase{TunnelEndPointRouterIP: routerStruct, Err: err}
}

func GetAll() tunnelEndPointRouterIP.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return tunnelEndPointRouterIP.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return tunnelEndPointRouterIP.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var routers []core.TunnelEndPointRouterIP
	err = db.Preload("TunnelEndPointRouter").Find(&routers).Error
	return tunnelEndPointRouterIP.ResultDatabase{TunnelEndPointRouterIP: routers, Err: err}
}
