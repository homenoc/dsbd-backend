package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/jpnicTech"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(jpnic *core.JPNICTech) (*core.JPNICTech, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnic, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&jpnic).Error
	return jpnic, err
}

func Delete(jpnic *core.JPNICTech) error {
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

	return db.Delete(jpnic).Error
}

func Update(base int, jpnic core.JPNICTech) error {
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

	if base == jpnicTech.UpdateAll {
		err = db.Model(&core.JPNICTech{Model: gorm.Model{ID: jpnic.ID}}).Updates(jpnic).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return err
}

func Get(base int, data *core.JPNICTech) jpnicTech.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicTech.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return jpnicTech.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var jpnicTechStruct []core.JPNICTech

	if base == jpnicTech.ID { //ID
		err = db.First(&jpnicTechStruct, data.ID).Error
	} else {
		log.Println("base select error")
		return jpnicTech.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return jpnicTech.ResultDatabase{Tech: jpnicTechStruct, Err: err}
}

func GetAll() jpnicTech.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicTech.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return jpnicTech.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var networks []core.JPNICTech
	err = db.Find(&networks).Error
	return jpnicTech.ResultDatabase{Tech: networks, Err: err}
}
