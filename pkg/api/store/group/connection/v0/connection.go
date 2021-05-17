package v0

import (
	"fmt"
	core "github.com/homenoc/dsbd-backend/pkg/api/core"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(connection *core.Connection) (*core.Connection, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&connection).Error
	return connection, err
}

func Delete(connection *core.Connection) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(connection).Error
}

func Update(base int, c core.Connection) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if connection.UpdateInfo == base {
		result = db.Model(&core.Connection{Model: gorm.Model{ID: c.ID}}).Update(core.Connection{
			NTTTemplateID: c.NTTTemplateID,
			NOC:           c.NOC,
			TermIP:        c.TermIP,
			Monitor:       c.Monitor,
		})
	} else if connection.UpdateServiceID == base {
		result = db.Model(&core.Connection{Model: gorm.Model{ID: c.ID}}).Update(core.Connection{ServiceID: c.ServiceID})
	} else if base == connection.UpdateAll {
		result = db.Model(&core.Connection{Model: gorm.Model{ID: c.ID}}).Update(c)
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *core.Connection) connection.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var connectionStruct []core.Connection

	if base == connection.ID { //ID
		err = db.Preload("ConnectionTemplate").
			Preload("NOC").
			Preload("BGPRouter").
			Preload("TunnelEndPointRouterIP").
			Preload("NTTTemplate").
			Preload("Service").
			Preload("Service.ServiceTemplate").
			Preload("Service.Group").
			First(&connectionStruct, data.ID).Error
	} else if base == connection.ServiceID {
		err = db.Where("service_id = ?", data.ServiceID).Find(&connectionStruct).Error
	} else {
		log.Println("base select error")
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return connection.ResultDatabase{Connection: connectionStruct, Err: err}
}

func GetAll() connection.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var connections []core.Connection
	err = db.Preload("ConnectionTemplate").
		Preload("NTTTemplate").
		Preload("NOC").
		Preload("BGPRouter").
		Preload("TunnelEndPointRouterIP").
		Preload("Service").
		Preload("Service.ServiceTemplate").
		Preload("Service.Group").
		Find(&connections).Error
	return connection.ResultDatabase{Connection: connections, Err: err}

}
