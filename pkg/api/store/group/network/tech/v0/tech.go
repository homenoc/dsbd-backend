package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/tech"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(network *tech.Tech) (*tech.Tech, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return network, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&network).Error
	return network, err
}

func Delete(network *tech.Tech) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(network).Error
}

func Update(base int, u tech.Tech) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if base == tech.UpdateAll {
		err = db.Model(&tech.Tech{Model: gorm.Model{ID: u.ID}}).Update(tech.Tech{
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

func Get(base int, data *tech.Tech) tech.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return tech.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networkStruct []tech.Tech

	if base == tech.ID { //ID
		err = db.First(&networkStruct, data.ID).Error
	} else if base == tech.UserId { //Name
		err = db.Where("user_id = ?", data.UserID).Find(&networkStruct).Error
	} else if base == tech.NetworkID { //Name
		err = db.Where("network_id = ?", data.NetworkID).Find(&networkStruct).Error
	} else if base == tech.NetworkAndUserId {
		err = db.Where("network_id = ? AND user_id = ?", data.NetworkID, data.UserID).Find(&networkStruct).Error
	} else {
		log.Println("base select error")
		return tech.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return tech.ResultDatabase{Tech: networkStruct, Err: err}
}

func GetAll() tech.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return tech.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networks []tech.Tech
	err = db.Find(&networks).Error
	return tech.ResultDatabase{Tech: networks, Err: err}
}
