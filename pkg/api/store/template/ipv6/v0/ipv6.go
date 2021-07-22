package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	ipv6 "github.com/homenoc/dsbd-backend/pkg/api/core/template/ipv6"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(ipv4 *core.IPv6Template) (*core.IPv6Template, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv6 error")
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

func Delete(ipv4 *core.IPv6Template) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv6 error")
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

func Update(base int, c core.IPv6Template) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv6 error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = nil

	if base == ipv6.UpdateAll {
		err = db.Model(&core.IPv6Template{Model: gorm.Model{ID: c.ID}}).Updates(c).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return err
}

func Get(base int, data *core.IPv6Template) ipv6.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv6 error")
		return ipv6.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ipv6.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var ipv6Struct []core.IPv6Template

	if base == ipv6.ID { //ID
		err = db.First(&ipv6Struct, data.ID).Error
	} else if base == ipv6.Subnet { //Subnet
		err = db.Where("subnet = ?", data.Subnet).Find(&ipv6Struct).Error
	} else {
		log.Println("base select error")
		return ipv6.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return ipv6.ResultDatabase{IPv6: ipv6Struct, Err: err}
}

func GetAll() ipv6.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv6 error")
		return ipv6.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ipv6.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var ipv6s []core.IPv6Template
	err = db.Find(&ipv6s).Error
	return ipv6.ResultDatabase{IPv6: ipv6s, Err: err}

}
