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

func Create(ipv4 *core.IPv4RouteTemplate) (*core.IPv4RouteTemplate, error) {
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

func Delete(ipv4 *core.IPv4RouteTemplate) error {
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

func Update(base int, c core.IPv4RouteTemplate) error {
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
		err = db.Model(&core.IPv4RouteTemplate{Model: gorm.Model{ID: c.ID}}).Updates(c).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return err
}

func Get(id uint) ([]core.IPv4RouteTemplate, error) {
	var ipv4s []core.IPv4RouteTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return []core.IPv4RouteTemplate{}, err
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ipv4s, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.First(&ipv4s, id).Error
	return ipv4s, err
}

func GetAll() ([]core.IPv4RouteTemplate, error) {
	var ipv4s []core.IPv4RouteTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return []core.IPv4RouteTemplate{}, err
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ipv4s, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Find(&ipv4s).Error
	return ipv4s, err

}
