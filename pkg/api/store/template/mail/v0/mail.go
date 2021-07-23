package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"log"
	"time"
)

func Create(mail *core.MailTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&mail).Error
	return err
}

func Delete(mail *core.MailTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database mail error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	return db.Delete(mail).Error
}

func Update(base int, c core.MailTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database mail error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = nil

	return err
}

func Get(data *core.MailTemplate) error {
	//var mailStruct []core.MailTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Find(&data).Error

	return nil
}

func GetAll() ([]core.MailTemplate, error) {
	var mailStruct []core.MailTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return mailStruct, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return mailStruct, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Find(&mailStruct).Error
	return mailStruct, nil

}
