package v0

import (
	"fmt"
	core "github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouterIP"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(router *core.TunnelEndPointRouterIP) (*core.TunnelEndPointRouterIP, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return router, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&router).Error
	return router, err
}

func Delete(router *core.TunnelEndPointRouterIP) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(router).Error
}

func Update(base int, data core.TunnelEndPointRouterIP) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if tunnelEndPointRouterIP.UpdateAll == base {
		result = db.Model(&core.TunnelEndPointRouterIP{Model: gorm.Model{ID: data.ID}}).Update(core.TunnelEndPointRouterIP{
			IP:      data.IP,
			Comment: data.Comment,
			Enable:  data.Enable,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *core.TunnelEndPointRouterIP) tunnelEndPointRouterIP.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return tunnelEndPointRouterIP.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var routerStruct []core.TunnelEndPointRouterIP

	if base == tunnelEndPointRouterIP.ID { //ID
		err = db.First(&routerStruct, data.ID).Error
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
	defer db.Close()

	var routers []core.TunnelEndPointRouterIP
	err = db.Find(&routers).Error
	return tunnelEndPointRouterIP.ResultDatabase{TunnelEndPointRouterIP: routers, Err: err}
}
