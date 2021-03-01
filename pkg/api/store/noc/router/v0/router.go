package v0

import (
	"fmt"
	router "github.com/homenoc/dsbd-backend/pkg/api/core/noc/router"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(router *router.Router) (*router.Router, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return router, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&router).Error
	return router, err
}

func Delete(router *router.Router) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(router).Error
}

func Update(base int, data router.Router) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if router.UpdateAll == base {
		result = db.Model(&router.Router{Model: gorm.Model{ID: data.ID}}).Update(router.Router{
			NOC:      data.NOC,
			HostName: data.HostName,
			Address:  data.Address,
			Enable:   data.Enable,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *router.Router) router.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return router.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var routerStruct []router.Router

	if base == router.ID { //ID
		err = db.First(&routerStruct, data.ID).Error
	} else if base == router.NOC { //UserID
		err = db.Where("noc = ?", data.NOC).Find(&routerStruct).Error
	} else if base == router.Address { //UserID
		err = db.Where("address = ?", data.Address).Find(&routerStruct).Error
	} else if base == router.Enable { //GroupID
		err = db.Where("enable = ?", data.Enable).Find(&routerStruct).Error
	} else {
		log.Println("base select error")
		return router.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return router.ResultDatabase{Router: routerStruct, Err: err}
}

func GetAll() router.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return router.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var routers []router.Router
	err = db.Find(&routers).Error
	return router.ResultDatabase{Router: routers, Err: err}
}
