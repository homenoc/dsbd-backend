package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnic"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(network *jpnic.Jpnic) (*jpnic.Jpnic, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return network, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&network).Error
	return network, err
}

func Delete(network *jpnic.Jpnic) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(network).Error
}

func Update(base int, u jpnic.Jpnic) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if base == jpnic.UpdateAll {
		err = db.Model(&jpnic.Jpnic{Model: gorm.Model{ID: u.ID}}).Update(jpnic.Jpnic{
			NetworkId: u.NetworkId, UserId: u.UserId, Lock: u.Lock}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *jpnic.Jpnic) jpnic.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnic.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networkStruct []jpnic.Jpnic

	if base == jpnic.ID { //ID
		err = db.First(&networkStruct, data.ID).Error
	} else if base == jpnic.UserId { //Name
		err = db.Where("user_id = ?", data.UserId).Find(&networkStruct).Error
	} else if base == jpnic.NetworkId { //Name
		err = db.Where("network_id = ?", data.NetworkId).Find(&networkStruct).Error
	} else if base == jpnic.NetworkAndUserId {
		err = db.Where("network_id = ? AND user_id = ?", data.NetworkId, data.UserId).Find(&networkStruct).Error
	} else {
		log.Println("base select error")
		return jpnic.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return jpnic.ResultDatabase{Jpnic: networkStruct, Err: err}
}

func GetAll() jpnic.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnic.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networks []jpnic.Jpnic
	err = db.Find(&networks).Error
	return jpnic.ResultDatabase{Jpnic: networks, Err: err}
}
