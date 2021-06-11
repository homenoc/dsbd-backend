package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func JoinJPNICByAdmin(serviceID uint, input core.JPNICAdmin) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Model(&core.Service{Model: gorm.Model{ID: serviceID}}).
		Association("JPNICAdmin").
		Append(input).Error
}

func DeleteJPNICByAdmin(id uint) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(core.JPNICAdmin{Model: gorm.Model{ID: id}}).Error
}

func UpdateJPNICByAdmin(input core.JPNICAdmin) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Model(&core.JPNICAdmin{Model: gorm.Model{ID: input.ID}}).Update(input).Error
}
