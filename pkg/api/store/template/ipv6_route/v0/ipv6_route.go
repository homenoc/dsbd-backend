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

func Create(ipv6 *core.IPv6RouteTemplate) (*core.IPv6RouteTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv6 error")
		return ipv6, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&ipv6).Error
	return ipv6, err
}

func Delete(ipv4 *core.IPv6RouteTemplate) error {
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

func Update(base int, c core.IPv6RouteTemplate) error {
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
		err = db.Model(&core.IPv6RouteTemplate{Model: gorm.Model{ID: c.ID}}).Updates(c).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return err
}

func Get(id uint) ([]core.IPv6RouteTemplate, error) {
	var ipv6s []core.IPv6RouteTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return []core.IPv6RouteTemplate{}, err
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ipv6s, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.First(&ipv6s, id).Error
	return ipv6s, err
}

func GetAll() ([]core.IPv6RouteTemplate, error) {
	var ipv6s []core.IPv6RouteTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv6 error")
		return []core.IPv6RouteTemplate{}, err
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ipv6s, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Find(&ipv6s).Error
	return ipv6s, err

}
