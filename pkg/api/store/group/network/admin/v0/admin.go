package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/admin"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(network *admin.Admin) (*admin.Admin, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return network, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&network).Error
	return network, err
}

func Delete(network *admin.Admin) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(network).Error
}

func Update(base int, u admin.Admin) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if base == admin.UpdateAll {
		err = db.Model(&admin.Admin{Model: gorm.Model{ID: u.ID}}).Update(admin.Admin{
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

func Get(base int, data *admin.Admin) admin.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return admin.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networkStruct []admin.Admin

	if base == admin.ID { //ID
		err = db.First(&networkStruct, data.ID).Error
	} else if base == admin.UserId { //Name
		err = db.Where("user_id = ?", data.UserID).Find(&networkStruct).Error
	} else if base == admin.NetworkId { //Name
		err = db.Where("network_id = ?", data.NetworkID).Find(&networkStruct).Error
	} else if base == admin.NetworkAndUserId {
		err = db.Where("network_id = ? AND user_id = ?", data.NetworkID, data.UserID).Find(&networkStruct).Error
	} else {
		log.Println("base select error")
		return admin.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return admin.ResultDatabase{Admins: networkStruct, Err: err}
}

func GetAll() admin.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return admin.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networks []admin.Admin
	err = db.Find(&networks).Error
	return admin.ResultDatabase{Admins: networks, Err: err}
}
