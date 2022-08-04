package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/payment"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(input *core.Payment) error {
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

	return db.Create(&input).Error
}

func Delete(input *core.Payment) error {
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

	return db.Delete(&input).Error
}

func Update(base int, input *core.Payment) error {
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

	if payment.UpdatePaid == base {
		err = db.Model(&core.Payment{PaymentIntentID: input.PaymentIntentID}).Where("payment_intent_id", input.PaymentIntentID).Update("paid", input.Paid).Error
	} else if payment.UpdateAll == base {
		err = db.Model(&core.Payment{Model: gorm.Model{ID: input.ID}}).Updates(&input).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return err
}

func Get(base uint, input core.Payment) ([]core.Payment, error) {
	var payments []core.Payment

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return payments, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return payments, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	if base == payment.ID { //ID
		err = db.First(&payments, input.ID).Error
	} else if base == payment.PaymentIntentID { //PaymentIntentID
		err = db.Where("payment_intent_id = ?", input.PaymentIntentID).Find(&payments).Error
	} else {
		log.Println("base select error")
		return payments, fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return payments, err
}

func GetAll() ([]core.Payment, error) {
	var payments []core.Payment

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return payments, err
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return payments, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Preload("User").Preload("Group").Find(&payments).Error
	return payments, err
}
