package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/template/service"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(connection *core.ServiceTemplate) (*core.ServiceTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&connection).Error
	return connection, err
}

func Delete(connection *core.ServiceTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(connection).Error
}

func Update(base int, c core.ServiceTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if base == service.UpdateAll {
		result = db.Model(&core.ServiceTemplate{Model: gorm.Model{ID: c.ID}}).Update(core.ServiceTemplate{
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

func Get(base int, data *core.ServiceTemplate) service.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return service.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var connectionStruct []core.ServiceTemplate

	if base == service.ID { //ID
		err = db.First(&connectionStruct, data.ID).Error
	} else {
		log.Println("base select error")
		return service.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}

	return service.ResultDatabase{Services: connectionStruct, Err: err}
}

func GetAll() service.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return service.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var connections []core.ServiceTemplate
	err = db.Find(&connections).Error

	return service.ResultDatabase{Services: connections, Err: err}
}
