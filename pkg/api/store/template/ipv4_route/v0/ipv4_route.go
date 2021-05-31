package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	ipv4 "github.com/homenoc/dsbd-backend/pkg/api/core/template/ipv4"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(ipv4 *core.IPv4RouteTemplate) (*core.IPv4RouteTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return ipv4, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&ipv4).Error
	return ipv4, err
}

func Delete(ipv4 *core.IPv4RouteTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(ipv4).Error
}

func Update(base int, c core.IPv4RouteTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if base == ipv4.UpdateAll {
		result = db.Model(&core.IPv4RouteTemplate{Model: gorm.Model{ID: c.ID}}).Update(c)
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return result.Error
}

func GetAll() ([]core.IPv4RouteTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ipv4 error")
		return []core.IPv4RouteTemplate{}, err
	}
	defer db.Close()

	var ipv4s []core.IPv4RouteTemplate
	err = db.Find(&ipv4s).Error
	return ipv4s, err

}
