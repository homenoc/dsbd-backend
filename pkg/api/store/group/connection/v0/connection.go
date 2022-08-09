package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(connection *core.Connection) (*core.Connection, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&connection).Error
	return connection, err
}

func Delete(connection *core.Connection) error {
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

	return db.Delete(connection).Error
}

func Update(base int, c core.Connection) error {
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

	if connection.UpdateInfo == base {
		err = db.Model(&core.Connection{Model: gorm.Model{ID: c.ID}}).Updates(core.Connection{
			NTT:     c.NTT,
			NOC:     c.NOC,
			TermIP:  c.TermIP,
			Monitor: c.Monitor,
		}).Error
	} else if connection.UpdateServiceID == base {
		err = db.Model(&core.Connection{Model: gorm.Model{ID: c.ID}}).Updates(core.Connection{ServiceID: c.ServiceID}).Error
	} else if base == connection.UpdateAll {
		err = db.Model(&core.Connection{Model: gorm.Model{ID: c.ID}}).Updates(c).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return err
}

func Get(base int, data *core.Connection) connection.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var connectionStruct []core.Connection

	if base == connection.ID { //ID
		err = db.Preload("NOC").
			Preload("BGPRouter").
			Preload("TunnelEndPointRouterIP").
			Preload("Service").
			Preload("Service.Group").
			First(&connectionStruct, data.ID).Error
	} else if base == connection.ServiceID {
		err = db.Where("service_id = ?", data.ServiceID).Find(&connectionStruct).Error
	} else if base == connection.NOCID {
		err = db.Preload("NOC").
			Preload("BGPRouter").
			Preload("TunnelEndPointRouterIP").
			Preload("Service").
			Preload("Service.Group").
			Preload("Service.Group.User").
			Where("noc_id = ?", data.NOCID).
			Find(&connectionStruct).Error
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
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var connections []core.Connection
	err = db.Preload("NOC").
		Preload("BGPRouter").
		Preload("BGPRouter.NOC").
		Preload("TunnelEndPointRouterIP").
		Preload("Service").
		Preload("Service.Group").
		Find(&connections).Error
	return connection.ResultDatabase{Connection: connections, Err: err}

}
