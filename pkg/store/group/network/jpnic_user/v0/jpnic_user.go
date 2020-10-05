package v0

import (
	"fmt"
	jpnicUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnic_user"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(network *jpnicUser.JPNICUser) (*jpnicUser.JPNICUser, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return network, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&network).Error
	return network, err
}

func Delete(network *jpnicUser.JPNICUser) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(network).Error
}

func Update(base int, u jpnicUser.JPNICUser) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if jpnicUser.UpdateInfo == base {
		result = db.Model(&jpnicUser.JPNICUser{Model: gorm.Model{ID: u.ID}}).Update(jpnicUser.JPNICUser{
			NameJa: u.NameJa, Name: u.Name, OrgJa: u.OrgJa, Org: u.Org, PostCode: u.PostCode, AddressJa: u.AddressJa,
			Address: u.Address, DeptJa: u.DeptJa, Dept: u.Dept, PosJa: u.PosJa, Pos: u.Pos,
			Mail: u.Mail, Tel: u.Tel, Fax: u.Fax})
	} else if jpnicUser.UpdateGID == base {
		result = db.Model(&jpnicUser.JPNICUser{Model: gorm.Model{ID: u.ID}}).Update(jpnicUser.JPNICUser{GroupID: u.GroupID})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *jpnicUser.JPNICUser) jpnicUser.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicUser.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networkStruct []jpnicUser.JPNICUser

	if base == jpnicUser.ID { //ID
		err = db.First(&networkStruct, jpnicUser.ID).Error
	} else if base == jpnicUser.Name { //Name
		err = db.Where("name = ?", jpnicUser.Name).Find(&networkStruct).Error
	} else if base == jpnicUser.Mail { //Name
		err = db.Where("mail = ?", jpnicUser.Mail).Find(&networkStruct).Error
	} else if base == jpnicUser.GID {
		err = db.Where("group_id = ?", jpnicUser.GID).Find(&networkStruct).Error
	} else {
		log.Println("base select error")
		return jpnicUser.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return jpnicUser.ResultDatabase{JPNICUser: networkStruct, Err: err}
}

func GetAll() jpnicUser.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicUser.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networks []jpnicUser.JPNICUser
	err = db.Find(&networks).Error
	return jpnicUser.ResultDatabase{JPNICUser: networks, Err: err}
}
