package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	ipv4 "github.com/homenoc/dsbd-backend/pkg/api/core/template/ipv4"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(ipv4 *core.IPv4Template) (*core.IPv4Template, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return ipv4, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&ipv4).Error
	return ipv4, err
}

func Delete(ipv4 *core.IPv4Template) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	return db.Delete(ipv4).Error
}

func Update(base int, c core.IPv4Template) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = nil

	if base == ipv4.UpdateAll {
		err = db.Model(&core.IPv4Template{Model: gorm.Model{ID: c.ID}}).Updates(c).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return err
}

func Get(base int, data *core.IPv4Template) ipv4.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return ipv4.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ipv4.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var ipv4Struct []core.IPv4Template

	if base == ipv4.ID { //ID
		err = db.First(&ipv4Struct, data.ID).Error
	} else if base == ipv4.Subnet { //Subnet
		err = db.Where("subnet = ?", data.Subnet).Find(&ipv4Struct).Error
	} else {
		log.Println("base select error")
		return ipv4.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return ipv4.ResultDatabase{IPv4: ipv4Struct, Err: err}
}

func GetAll() ipv4.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return ipv4.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ipv4.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var ipv4s []core.IPv4Template
	err = db.Find(&ipv4s).Error
	return ipv4.ResultDatabase{IPv4: ipv4s, Err: err}

}
