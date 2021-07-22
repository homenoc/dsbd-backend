package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(connection *core.PaymentDonateTemplate) (*core.PaymentDonateTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&connection).Error
	return connection, err
}

func Delete(connection *core.PaymentDonateTemplate) error {
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

	return db.Delete(connection).Error
}

func Update(c core.PaymentDonateTemplate) error {
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

	err = db.Model(&core.PaymentDonateTemplate{Model: gorm.Model{ID: c.ID}}).Updates(core.PaymentDonateTemplate{
		Name:    c.Name,
		Fee:     c.Fee,
		Comment: c.Comment,
	}).Error

	return err
}

func Get(id uint) (core.PaymentDonateTemplate, error) {
	var payment core.PaymentDonateTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return payment, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return payment, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.First(&payment, id).Error

	return payment, err
}

func GetAll() ([]core.PaymentDonateTemplate, error) {
	var payments []core.PaymentDonateTemplate

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

	err = db.Find(&payments).Error

	return payments, err
}
