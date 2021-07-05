package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(connection *core.PaymentMembershipTemplate) (*core.PaymentMembershipTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return connection, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&connection).Error
	return connection, err
}

func Delete(connection *core.PaymentMembershipTemplate) error {
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

func Update(c core.PaymentMembershipTemplate) error {
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

	err = db.Model(&core.PaymentMembershipTemplate{Model: gorm.Model{ID: c.ID}}).Updates(core.PaymentMembershipTemplate{
		Title:   c.Title,
		Plan:    c.Plan,
		Comment: c.Comment,
	}).Error

	return err
}

func Get(id uint) (core.PaymentMembershipTemplate, error) {
	var payment core.PaymentMembershipTemplate

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

func GetAll() ([]core.PaymentMembershipTemplate, error) {
	var payments []core.PaymentMembershipTemplate

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
