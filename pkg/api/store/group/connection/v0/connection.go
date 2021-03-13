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
		result = db.Model(&core.Connection{Model: gorm.Model{ID: c.ID}}).Update(core.Connection{
			ServiceID:                c.ServiceID,
			BGPRouterID:              c.BGPRouterID,
			TunnelEndPointRouterIPID: c.TunnelEndPointRouterIPID,
			NTTTemplateID:            c.NTTTemplateID,
			NOC:                      c.NOC,
			TermIP:                   c.TermIP,
			Monitor:                  c.Monitor,
			LinkV4Our:                c.LinkV4Our,
			LinkV4Your:               c.LinkV4Your,
			LinkV6Our:                c.LinkV6Our,
			LinkV6Your:               c.LinkV6Your,
			Open:                     c.Open,
			Lock:                     c.Lock,
		})
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
		err = db.First(&connectionStruct, data.ID).Error
	} else if base == connection.ServiceID {
		err = db.Where("group_id = ?", data.ServiceID).Find(&connectionStruct).Error
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
	err = db.Find(&connections).Error
	return connection.ResultDatabase{Connection: connections, Err: err}

}
