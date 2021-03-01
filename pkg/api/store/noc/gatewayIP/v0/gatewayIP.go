package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/gatewayIP"
	router "github.com/homenoc/dsbd-backend/pkg/api/core/noc/router"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(router *gatewayIP.GatewayIP) (*gatewayIP.GatewayIP, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return router, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&router).Error
	return router, err
}

func Delete(router *gatewayIP.GatewayIP) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(router).Error
}

func Update(base int, data gatewayIP.GatewayIP) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if router.UpdateAll == base {
		result = db.Model(&gatewayIP.GatewayIP{Model: gorm.Model{ID: data.ID}}).Update(gatewayIP.GatewayIP{
			GatewayID: data.GatewayID,
			IP:        data.IP,
			Comment:   data.Comment,
			Enable:    data.Enable,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *gatewayIP.GatewayIP) gatewayIP.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return gatewayIP.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var routerStruct []gatewayIP.GatewayIP

	if base == router.ID { //ID
		err = db.First(&routerStruct, data.ID).Error
	} else if base == router.Enable { //GroupID
		err = db.Where("enable = ?", data.Enable).Find(&routerStruct).Error
	} else {
		log.Println("base select error")
		return gatewayIP.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return gatewayIP.ResultDatabase{GatewayIP: routerStruct, Err: err}
}

func GetAll() gatewayIP.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return gatewayIP.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var routers []gatewayIP.GatewayIP
	err = db.Find(&routers).Error
	return gatewayIP.ResultDatabase{GatewayIP: routers, Err: err}
}
