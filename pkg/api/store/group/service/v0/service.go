package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(service *core.Service) (*core.Service, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return service, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&service).Error
	return service, err
}

func Delete(service *core.Service) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(service).Error
}

func Update(base int, c core.Service) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if service.UpdateData == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Update(core.Service{
			Org:       c.Org,
			OrgEn:     c.OrgEn,
			Postcode:  c.Postcode,
			Address:   c.Address,
			AddressEn: c.AddressEn,
			RouteV4:   c.RouteV4,
			RouteV6:   c.RouteV6,
			ASN:       c.ASN,
		})
	} else if service.UpdateRoute == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Update(core.Service{
			RouteV4: c.RouteV4, RouteV6: c.RouteV6})
	} else if service.UpdateGID == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Update(core.Service{GroupID: c.GroupID})
	} else if service.UpdateStatus == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Update(core.Service{AddAllow: c.AddAllow})
	} else if service.UpdateAll == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Update(core.Service{
			GroupID:   c.GroupID,
			Org:       c.Org,
			OrgEn:     c.Org,
			Postcode:  c.Postcode,
			Address:   c.Address,
			AddressEn: c.AddressEn,
			ASN:       c.ASN,
			V4Name:    c.V4Name,
			V6Name:    c.V6Name,
			RouteV4:   c.RouteV4,
			RouteV6:   c.RouteV6,
			IP:        c.IP,
			Open:      c.Open,
			Lock:      c.Lock,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func JoinJPNICTech(serviceID, jpnicTechID uint) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	db.Model(&core.Service{Model: gorm.Model{ID: serviceID}}).Association("JPNICTech").
		Append(&core.JPNICTech{Model: gorm.Model{ID: jpnicTechID}})

	return nil
}

func DeleteJPNICTech(serviceID, jpnicTechID uint) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	db.Model(&core.Service{Model: gorm.Model{ID: serviceID}}).Association("JPNICTech").
		Delete(&core.JPNICTech{Model: gorm.Model{ID: jpnicTechID}})

	return nil
}

func Get(base int, data *core.Service) service.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return service.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var serviceStruct []core.Service

	if base == service.ID { //ID
		err = db.Preload("ServiceTemplate").
			Preload("IP").
			Preload("IP.Plan").
			Preload("Connection").
			Preload("Connection.ConnectionTemplate").
			Preload("Connection.NOC").
			Preload("Connection.BGPRouter").
			Preload("Connection.TunnelEndPointRouterIP").
			Preload("ServiceTemplate").
			Preload("JPNICAdmin").
			Preload("JPNICTech").
			First(&serviceStruct, data.ID).Error
	} else if base == service.Org { //Mail
		err = db.Preload("ServiceTemplate").
			Preload("IP").
			Preload("Connection").
			Preload("JPNICAdmin").
			Preload("JPNICTech").
			Where("org = ?", data.Org).Find(&serviceStruct).Error
	} else if base == service.GID {
		err = db.Preload("ServiceTemplate").
			Preload("IP").
			Preload("Connection").
			Preload("JPNICAdmin").
			Preload("JPNICTech").
			Where("group_id = ?", data.GroupID).Find(&serviceStruct).Error
	} else if base == service.GIDAndAddAllow {
		err = db.Preload("ServiceTemplate").
			Where("group_id = ? AND add_allow = ?", data.GroupID, true).Find(&serviceStruct).Error
	} else if base == service.SearchNewNumber {
		err = db.Where("group_id = ?", data.GroupID).Find(&serviceStruct).Error
	} else if base == service.Open {
		err = db.Where("group_id = ? AND open = ?", data.GroupID, true).
			Preload("IP", "open = ?", true).
			Preload("Connection", "open = ?", true).
			Preload("Connection.ConnectionTemplate").
			Preload("Connection.NOC").
			Preload("Connection.BGPRouter").
			Preload("Connection.TunnelEndPointRouterIP").
			Preload("ServiceTemplate").
			Preload("JPNICAdmin").
			Preload("JPNICTech").
			Find(&serviceStruct).Error
	} else {
		log.Println("base select error")
		return service.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return service.ResultDatabase{Err: err, Service: serviceStruct}
}

func GetAll() service.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return service.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var services []core.Service
	err = db.
		Preload("ServiceTemplate").
		Preload("IP").
		Preload("Connection").
		Preload("JPNICAdmin").
		Preload("JPNICTech").
		Find(&services).Error
	return service.ResultDatabase{Err: err, Service: services}
}
