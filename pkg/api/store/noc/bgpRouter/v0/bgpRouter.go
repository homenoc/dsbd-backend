package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	router "github.com/homenoc/dsbd-backend/pkg/api/core/noc/bgpRouter"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(router *core.BGPRouter) (*core.BGPRouter, error) {
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

func Delete(router *core.BGPRouter) error {
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

func Update(base int, data core.BGPRouter) error {
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
		err = db.Model(&core.BGPRouter{Model: gorm.Model{ID: data.ID}}).Updates(core.BGPRouter{
			HostName: data.HostName,
			Address:  data.Address,
			Enable:   data.Enable,
		}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return err
}

func Get(base int, data *core.BGPRouter) router.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return router.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return router.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var routerStruct []core.BGPRouter

	if base == router.ID { //ID
		err = db.First(&routerStruct, data.ID).Error
	} else if base == router.Address { //Address
		err = db.Where("address = ?", data.Address).Find(&routerStruct).Error
	} else if base == router.Enable { //GroupID
		err = db.Where("enable = ?", data.Enable).Find(&routerStruct).Error
	} else {
		log.Println("base select error")
		return router.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return router.ResultDatabase{BGPRouter: routerStruct, Err: err}
}

func GetAll() router.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return router.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return router.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var routers []core.BGPRouter
	err = db.Find(&routers).Error
	return router.ResultDatabase{BGPRouter: routers, Err: err}
}
