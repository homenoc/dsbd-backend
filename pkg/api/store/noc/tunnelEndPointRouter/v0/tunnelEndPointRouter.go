package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouter"
	router "github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouterIP"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(router *core.TunnelEndPointRouter) (*core.TunnelEndPointRouter, error) {
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

func Delete(router *core.TunnelEndPointRouter) error {
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

func Update(base int, data core.TunnelEndPointRouter) error {
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

	if router.UpdateAll == base {
		err = db.Model(&core.TunnelEndPointRouter{Model: gorm.Model{ID: data.ID}}).Updates(core.TunnelEndPointRouter{
			NOCID:    data.NOCID,
			HostName: data.HostName,
			Capacity: data.Capacity,
			Comment:  data.Comment,
			Enable:   data.Enable,
		}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return err
}

func Get(base int, data *core.TunnelEndPointRouter) tunnelEndPointRouter.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return tunnelEndPointRouter.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return tunnelEndPointRouter.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var routerStruct []core.TunnelEndPointRouter

	if base == router.ID { //ID
		err = db.First(&routerStruct, data.ID).Error
	} else if base == router.Enable { //GroupID
		err = db.Where("enable = ?", data.Enable).Find(&routerStruct).Error
	} else {
		log.Println("base select error")
		return tunnelEndPointRouter.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return tunnelEndPointRouter.ResultDatabase{TunnelEndPointRouter: routerStruct, Err: err}
}

func GetAll() tunnelEndPointRouter.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return tunnelEndPointRouter.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return tunnelEndPointRouter.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var routers []core.TunnelEndPointRouter
	err = db.Find(&routers).Error
	return tunnelEndPointRouter.ResultDatabase{TunnelEndPointRouter: routers, Err: err}
}
