package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/gateway"
	router "github.com/homenoc/dsbd-backend/pkg/api/core/noc/router"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(router *gateway.Gateway) (*gateway.Gateway, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return router, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&router).Error
	return router, err
}

func Delete(router *gateway.Gateway) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(router).Error
}

func Update(base int, data gateway.Gateway) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if router.UpdateAll == base {
		result = db.Model(&gateway.Gateway{Model: gorm.Model{ID: data.ID}}).Update(gateway.Gateway{
			NOCID:    data.NOCID,
			HostName: data.HostName,
			Capacity: data.Capacity,
			Comment:  data.Comment,
			Enable:   data.Enable,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *gateway.Gateway) gateway.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return gateway.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var routerStruct []gateway.Gateway

	if base == router.ID { //ID
		err = db.First(&routerStruct, data.ID).Error
	} else if base == router.Enable { //GroupID
		err = db.Where("enable = ?", data.Enable).Find(&routerStruct).Error
	} else {
		log.Println("base select error")
		return gateway.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return gateway.ResultDatabase{Gateway: routerStruct, Err: err}
}

func GetAll() gateway.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return gateway.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var routers []gateway.Gateway
	err = db.Find(&routers).Error
	return gateway.ResultDatabase{Gateway: routers, Err: err}
}
