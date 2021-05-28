package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func JoinPlan(ipID uint, input core.Plan) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Model(&core.IP{Model: gorm.Model{ID: ipID}}).
		Association("Plan").
		Append(&input).Error
}

func DeletePlan(id uint) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(core.Plan{Model: gorm.Model{ID: id}}).Error
}

func UpdatePlan(input core.Plan) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Model(&core.Plan{Model: gorm.Model{ID: input.ID}}).Update(input).Error
}

func GetPlan(data *core.Plan) (core.Plan, error) {
	var plans core.Plan

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return plans, err
	}
	defer db.Close()

	err = db.First(&plans, data.ID).Error

	return plans, err
}
