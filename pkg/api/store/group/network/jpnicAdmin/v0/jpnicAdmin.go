package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(network *jpnicAdmin.JpnicAdmin) (*jpnicAdmin.JpnicAdmin, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return network, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&network).Error
	return network, err
}

func Delete(network *jpnicAdmin.JpnicAdmin) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(network).Error
}

func Update(base int, u jpnicAdmin.JpnicAdmin) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if base == jpnicAdmin.UpdateAll {
		err = db.Model(&jpnicAdmin.JpnicAdmin{Model: gorm.Model{ID: u.ID}}).Update(jpnicAdmin.JpnicAdmin{
			NetworkID: u.NetworkID,
			UserID:    u.UserID,
			Lock:      u.Lock,
		}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *jpnicAdmin.JpnicAdmin) jpnicAdmin.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicAdmin.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networkStruct []jpnicAdmin.JpnicAdmin

	if base == jpnicAdmin.ID { //ID
		err = db.First(&networkStruct, data.ID).Error
	} else if base == jpnicAdmin.UserId { //Name
		err = db.Where("user_id = ?", data.UserID).Find(&networkStruct).Error
	} else if base == jpnicAdmin.NetworkId { //Name
		err = db.Where("network_id = ?", data.NetworkID).Find(&networkStruct).Error
	} else if base == jpnicAdmin.NetworkAndUserId {
		err = db.Where("network_id = ? AND user_id = ?", data.NetworkID, data.UserID).Find(&networkStruct).Error
	} else {
		log.Println("base select error")
		return jpnicAdmin.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return jpnicAdmin.ResultDatabase{Jpnic: networkStruct, Err: err}
}

func GetAll() jpnicAdmin.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicAdmin.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networks []jpnicAdmin.JpnicAdmin
	err = db.Find(&networks).Error
	return jpnicAdmin.ResultDatabase{Jpnic: networks, Err: err}
}
