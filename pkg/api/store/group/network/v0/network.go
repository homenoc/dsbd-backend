package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(network *network.Network) (*network.Network, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return network, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&network).Error
	return network, err
}

func Delete(network *network.Network) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(network).Error
}

func Update(base int, c network.Network) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if network.UpdatePlan == base {
		result = db.Model(&network.Network{Model: gorm.Model{ID: c.ID}}).Update(network.Network{Plan: c.Plan})
	} else if network.UpdateData == base {
		result = db.Model(&network.Network{Model: gorm.Model{ID: c.ID}}).Update(network.Network{
			Org: c.Org, OrgEn: c.OrgEn, Postcode: c.Postcode, Address: c.Address, AddressEn: c.AddressEn,
			RouteV4: c.RouteV4, RouteV6: c.RouteV6, PI: c.PI, ASN: c.ASN, V4: c.V4, V6: c.V6,
			V4Name: c.V4Name, V6Name: c.V6Name, Date: c.Date, Plan: c.Plan})
	} else if network.UpdateRoute == base {
		result = db.Model(&network.Network{Model: gorm.Model{ID: c.ID}}).Update(network.Network{
			RouteV4: c.RouteV4, RouteV6: c.RouteV6})
	} else if network.UpdateGID == base {
		result = db.Model(&network.Network{Model: gorm.Model{ID: c.ID}}).Update(network.Network{GroupID: c.GroupID})
	} else if network.UpdateAll == base {
		result = db.Model(&network.Network{Model: gorm.Model{ID: c.ID}}).Update(network.Network{
			GroupID: c.GroupID, Org: c.Org, OrgEn: c.Org, Postcode: c.Postcode, Address: c.Address, AddressEn: c.AddressEn,
			PI: c.PI, ASN: c.ASN, RouteV4: c.RouteV4, RouteV6: c.RouteV6, IP: c.IP, V4Name: c.V4Name, V6Name: c.V6Name,
			JPNICAdmin: c.JPNICAdmin, JPNICTech: c.JPNICTech, Open: c.Open, Lock: c.Lock})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *network.Network) network.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return network.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networkStruct []network.Network

	if base == network.ID { //ID
		err = db.Preload("IP").Preload("Connection").Preload("JPNICAdmin").Preload("JPNICTech").
			First(&networkStruct, data.ID).Error
	} else if base == network.Org { //Mail
		err = db.Preload("IP").Preload("Connection").Preload("JPNICAdmin").Preload("JPNICTech").
			Where("org = ?", data.Org).Find(&networkStruct).Error
	} else if base == network.GID {
		err = db.Preload("IP").Preload("Connection").Preload("JPNICAdmin").Preload("JPNICTech").
			Where("group_id = ?", data.GroupID).Find(&networkStruct).Error
	} else if base == network.Open {
		err = db.Preload("IP", "open = 1").Preload("Connection", "open = 1").Preload("JPNICAdmin").Preload("JPNICTech").
			Where("group_id = ? AND open = ?", data.GroupID, true).Find(&networkStruct).Error
	} else {
		log.Println("base select error")
		return network.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return network.ResultDatabase{Network: networkStruct, Err: err}
}

func GetAll() network.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return network.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var networks []network.Network
	err = db.Preload("IP").Preload("Connection").Preload("JPNICAdmin").Preload("JPNICTech").Find(&networks).Error
	return network.ResultDatabase{Network: networks, Err: err}
}
