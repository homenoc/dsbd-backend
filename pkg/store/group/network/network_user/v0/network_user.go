package v0

import (
	"fmt"
	networkUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/network_user"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(networkUser *networkUser.NetworkUser) (*networkUser.NetworkUser, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return networkUser, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&networkUser).Error
	return networkUser, err
}

func Delete(network *networkUser.NetworkUser) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(network).Error
}

func Update(base int, user networkUser.NetworkUser) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if base == networkUser.UpdateAll {
		err = db.Model(&networkUser.NetworkUser{Model: gorm.Model{ID: user.ID}}).Update(networkUser.NetworkUser{
			Type: user.Type, NetworkID: user.NetworkID, JPNICUserID: user.JPNICUserID}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *networkUser.NetworkUser) networkUser.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return networkUser.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networkStruct []networkUser.NetworkUser

	if base == networkUser.ID { //ID
		err = db.First(&networkStruct, data.ID).Error
	} else if base == networkUser.NetAndType { //Name
		err = db.Where("type = ? AND network_id = ?", data.Type, data.NetworkID).Find(&networkStruct).Error
	} else if base == networkUser.Network { //Name
		err = db.Where("network_id = ?", data.NetworkID).Find(&networkStruct).Error
	} else if base == networkUser.User {
		err = db.Where("jpnic_user_id = ?", data.JPNICUserID).Find(&networkStruct).Error
	} else if base == networkUser.Type {
		err = db.Where("type = ?", data.Type).Find(&networkStruct).Error
	} else if base == networkUser.Info {
		err = db.Where("type = ? AND network_id = ? AND jpnic_user_id = ?", data.Type, data.NetworkID,
			data.JPNICUserID).Find(&networkStruct).Error
	} else {
		log.Println("base select error")
		return networkUser.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return networkUser.ResultDatabase{NetworkUser: networkStruct, Err: err}
}

func GetAll() networkUser.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return networkUser.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networks []networkUser.NetworkUser
	err = db.Find(&networks).Error
	return networkUser.ResultDatabase{NetworkUser: networks, Err: err}
}
