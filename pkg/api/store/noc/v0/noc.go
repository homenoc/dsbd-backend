package v0

import (
	"fmt"
	core "github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(noc *core.NOC) (*core.NOC, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return noc, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&noc).Error
	return noc, err
}

func Delete(noc *core.NOC) error {
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

	return db.Delete(noc).Error
}

func Update(base int, data core.NOC) error {
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

	if noc.UpdateAll == base {
		err = db.Model(&core.NOC{Model: gorm.Model{ID: data.ID}}).Updates(core.NOC{
			Name:      data.Name,
			Location:  data.Location,
			Bandwidth: data.Bandwidth,
			Enable:    data.Enable,
			Comment:   data.Comment,
		}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return err
}

func Get(base int, data *core.NOC) noc.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return noc.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return noc.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var nocStruct []core.NOC

	if base == noc.ID { //ID
		err = db.First(&nocStruct, data.ID).Error
	} else if base == noc.Name { //UserID
		err = db.Where("name = ?", data.Name).Find(&nocStruct).Error
	} else if base == noc.Enable { //GroupID
		err = db.Where("enable = ?", data.Enable).Find(&nocStruct).Error
	} else {
		log.Println("base select error")
		return noc.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return noc.ResultDatabase{NOC: nocStruct, Err: err}
}

func GetAll() noc.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return noc.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return noc.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var nocs []core.NOC
	err = db.Preload("BGPRouter").
		Preload("TunnelEndPointRouter").
		Preload("TunnelEndPointRouter.TunnelEndPointRouterIP").
		Find(&nocs).Error
	return noc.ResultDatabase{NOC: nocs, Err: err}
}
