package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(plan *core.Plan) (*core.Plan, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return plan, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&plan).Error
	return plan, err
}

func Delete(plan *core.Plan) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(plan).Error
}

func Update(u core.Plan) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Model(&core.Plan{Model: gorm.Model{ID: u.ID}}).Update(u).Error

	return err
}
