package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(connection *core.PaymentDonateTemplate) (*core.PaymentDonateTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&connection).Error
	return connection, err
}

func Delete(connection *core.PaymentDonateTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(connection).Error
}

func Update(c core.PaymentDonateTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	result = db.Model(&core.PaymentDonateTemplate{Model: gorm.Model{ID: c.ID}}).Update(core.PaymentDonateTemplate{
		Name:    c.Name,
		Fee:     c.Fee,
		Comment: c.Comment,
	})

	return result.Error
}

func Get(id uint) (core.PaymentDonateTemplate, error) {
	var payment core.PaymentDonateTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return payment, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.First(&payment, id).Error

	return payment, err
}

func GetAll() ([]core.PaymentDonateTemplate, error) {
	var connections []core.PaymentDonateTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connections, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Find(&connections).Error

	return connections, err
}
