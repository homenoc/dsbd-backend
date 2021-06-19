package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/payment"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(input *core.Payment) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}

	defer db.Close()

	return db.Create(&input).Error
}

func Delete(input *core.Payment) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(&input).Error
}

func Update(base int, input *core.Payment) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if payment.UpdatePaid == base {
		result = db.Model(&core.Payment{PaymentIntentID: input.PaymentIntentID}).Update(core.Payment{Paid: input.Paid})
	} else if payment.UpdateAll == base {
		result = db.Model(&core.Payment{Model: gorm.Model{ID: input.ID}}).Update(&input)
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return result.Error
}

func Get(base uint, input core.Payment) ([]core.Payment, error) {
	var payments []core.Payment

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return payments, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

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
	defer db.Close()

	err = db.Preload("User").Preload("Group").Find(&payments).Error
	return payments, err
}
