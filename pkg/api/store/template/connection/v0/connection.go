package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/template/connection"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(connection *core.ConnectionTemplate) (*core.ConnectionTemplate, error) {
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

func Delete(connection *core.ConnectionTemplate) error {
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

func Update(base int, c core.ConnectionTemplate) error {
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

	if base == connection.UpdateAll {
		err = db.Model(&core.ConnectionTemplate{Model: gorm.Model{ID: c.ID}}).Updates(core.ConnectionTemplate{
			Hidden:  c.Hidden,
			Name:    c.Name,
			Comment: c.Comment,
		}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return err
}

func Get(base int, data *core.ConnectionTemplate) connection.ResultDatabase {
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
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var connections []core.ConnectionTemplate
	err = db.Find(&connections).Error
	return connection.ResultDatabase{Connections: connections, Err: err}

}
