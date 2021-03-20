package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/template/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(connection *core.ConnectionTemplate) (*core.ConnectionTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&connection).Error
	return connection, err
}

func Delete(connection *core.ConnectionTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(connection).Error
}

func Update(base int, c core.ConnectionTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if base == connection.UpdateAll {
		result = db.Model(&core.ConnectionTemplate{Model: gorm.Model{ID: c.ID}}).Update(core.ConnectionTemplate{
			Hidden:  c.Hidden,
			Name:    c.Name,
			Comment: c.Comment,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return result.Error
}

func Get(base int, data *core.ConnectionTemplate) connection.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var connectionStruct []core.ConnectionTemplate

	if base == connection.ID { //ID
		err = db.First(&connectionStruct, data.ID).Error
	} else {
		log.Println("base select error")
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return connection.ResultDatabase{Connections: connectionStruct, Err: err}
}

func GetAll() connection.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var connections []core.ConnectionTemplate
	err = db.Find(&connections).Error
	return connection.ResultDatabase{Connections: connections, Err: err}

}
