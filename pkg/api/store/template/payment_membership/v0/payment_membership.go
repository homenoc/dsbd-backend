package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(connection *core.PaymentMembershipTemplate) (*core.PaymentMembershipTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&connection).Error
	return connection, err
}

func Delete(connection *core.PaymentMembershipTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(connection).Error
}

func Update(c core.PaymentMembershipTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	result = db.Model(&core.PaymentMembershipTemplate{Model: gorm.Model{ID: c.ID}}).Update(core.PaymentMembershipTemplate{
		Title:   c.Title,
		Plan:    c.Plan,
		Comment: c.Comment,
	})

	return result.Error
}

func Get(id uint) (core.PaymentMembershipTemplate, error) {
	var payment core.PaymentMembershipTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return payment, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.First(&payment, id).Error

	return payment, err
}

func GetAll() ([]core.PaymentMembershipTemplate, error) {
	var payments []core.PaymentMembershipTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return payments, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Find(&payments).Error

	return payments, err
}
