package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(noc *noc.NOC) (*noc.NOC, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return noc, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&noc).Error
	return noc, err
}

func Delete(noc *noc.NOC) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(noc).Error
}

func Update(base int, data noc.NOC) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if noc.UpdateAll == base {
		result = db.Model(&noc.NOC{Model: gorm.Model{ID: data.ID}}).Update(noc.NOC{
			Name: data.Name, Location: data.Location, Bandwidth: data.Bandwidth,
			Enable: data.Enable, Comment: data.Comment})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *noc.NOC) noc.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return noc.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var nocStruct []noc.NOC

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
	defer db.Close()

	var nocs []noc.NOC
	err = db.Find(&nocs).Error
	return noc.ResultDatabase{NOC: nocs, Err: err}
}
