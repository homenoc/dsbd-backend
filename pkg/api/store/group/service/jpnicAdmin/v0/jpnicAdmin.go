package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/jpnicAdmin"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(network *core.JPNICAdmin) (*core.JPNICAdmin, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return network, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&network).Error
	return network, err
}

func Delete(network *core.JPNICAdmin) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(network).Error
}

func Update(base int, u core.JPNICAdmin) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if base == jpnicAdmin.UpdateAll {
		err = db.Model(&core.JPNICAdmin{Model: gorm.Model{ID: u.ID}}).Update(core.JPNICAdmin{
			Org:       u.Org,
			OrgEn:     u.OrgEn,
			PostCode:  u.PostCode,
			Address:   u.Address,
			AddressEn: u.AddressEn,
			Dept:      u.Dept,
			DeptEn:    u.DeptEn,
			Tel:       u.Tel,
			Fax:       u.Fax,
			Country:   u.Country,
		}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *core.JPNICAdmin) jpnicAdmin.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicAdmin.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var jpnicAdminStruct []core.JPNICAdmin

	if base == jpnicAdmin.ID { //ID
		err = db.First(&jpnicAdminStruct, data.ID).Error
	} else {
		log.Println("base select error")
		return jpnicAdmin.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return jpnicAdmin.ResultDatabase{Admins: jpnicAdminStruct, Err: err}
}

func GetAll() jpnicAdmin.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicAdmin.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networks []core.JPNICAdmin
	err = db.Find(&networks).Error
	return jpnicAdmin.ResultDatabase{Admins: networks, Err: err}
}
