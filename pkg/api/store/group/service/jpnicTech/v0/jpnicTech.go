package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/jpnicTech"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(network *core.JPNICTech) (*core.JPNICTech, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return network, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&network).Error
	return network, err
}

func Delete(network *core.JPNICTech) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(network).Error
}

func Update(base int, u core.JPNICTech) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if base == jpnicTech.UpdateAll {
		err = db.Model(&core.JPNICTech{Model: gorm.Model{ID: u.ID}}).Update(core.JPNICTech{
			Org:       u.Org,
			OrgEn:     u.OrgEn,
			PostCode:  u.PostCode,
			Address:   u.Address,
			AddressEn: u.AddressEn,
			Dept:      u.Dept,
			DeptEn:    u.DeptEn,
			Pos:       u.Pos,
			PosEn:     u.PosEn,
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

func Get(base int, data *core.JPNICTech) jpnicTech.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicTech.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networkStruct []core.JPNICTech

	if base == jpnicTech.ID { //ID
		err = db.First(&networkStruct, data.ID).Error
	} else {
		log.Println("base select error")
		return jpnicTech.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return jpnicTech.ResultDatabase{Tech: networkStruct, Err: err}
}

func GetAll() jpnicTech.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return jpnicTech.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networks []core.JPNICTech
	err = db.Find(&networks).Error
	return jpnicTech.ResultDatabase{Tech: networks, Err: err}
}
