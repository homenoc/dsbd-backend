package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicTech"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(network *jpnicTech.JpnicTech) (*jpnicTech.JpnicTech, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return network, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&network).Error
	return network, err
}

func Delete(network *jpnicTech.JpnicTech) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(network).Error
}

func Update(base int, u jpnicTech.JpnicTech) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if base == jpnicTech.UpdateAll {
		err = db.Model(&jpnicTech.JpnicTech{Model: gorm.Model{ID: u.ID}}).Update(jpnicTech.JpnicTech{
			NetworkId: u.NetworkId, UserId: u.UserId, Lock: u.Lock}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *jpnicTech.JpnicTech) jpnicTech.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicTech.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networkStruct []jpnicTech.JpnicTech

	if base == jpnicTech.ID { //ID
		err = db.First(&networkStruct, data.ID).Error
	} else if base == jpnicTech.UserId { //Name
		err = db.Where("user_id = ?", data.UserId).Find(&networkStruct).Error
	} else if base == jpnicTech.NetworkId { //Name
		err = db.Where("network_id = ?", data.NetworkId).Find(&networkStruct).Error
	} else if base == jpnicTech.NetworkAndUserId {
		err = db.Where("network_id = ? AND user_id = ?", data.NetworkId, data.UserId).Find(&networkStruct).Error
	} else {
		log.Println("base select error")
		return jpnicTech.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return jpnicTech.ResultDatabase{Jpnic: networkStruct, Err: err}
}

func GetAll() jpnicTech.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicTech.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networks []jpnicTech.JpnicTech
	err = db.Find(&networks).Error
	return jpnicTech.ResultDatabase{Jpnic: networks, Err: err}
}
