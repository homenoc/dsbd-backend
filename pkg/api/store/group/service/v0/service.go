package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(service *core.Service) (*core.Service, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return service, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&service).Error
	return service, err
}

func Delete(service *core.Service) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	return db.Delete(service).Error
}

func Update(base int, c core.Service) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = nil

	if service.UpdateData == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Updates(core.Service{
			Org:       c.Org,
			OrgEn:     c.OrgEn,
			PostCode:  c.PostCode,
			Address:   c.Address,
			AddressEn: c.AddressEn,
			ASN:       c.ASN,
		}).Error
	} else if service.UpdateGID == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Updates(core.Service{GroupID: c.GroupID}).Error
	} else if service.UpdateStatus == base {
	} else if service.UpdateAll == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Updates(c).Error
	} else if service.ReplaceIP == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("IP").Replace(c.IP[0])
	} else if service.AppendIP == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("IP").Replace(c.IP[0])
	} else if service.AppendJPNICAdmin == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("JPNICAdmin").Append(c.JPNICAdmin)
	} else if service.AppendJPNICTech == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("JPNICTech").Append(c.JPNICTech[0])
	} else if service.AppendConnection == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("IPv4").Replace(c.Connection[0])
	} else if service.DeleteIP == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("IP").Replace(c.IP[0])
	} else if service.DeleteJPNICAdmin == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("JPNICAdmin").Delete(c.JPNICAdmin)
	} else if service.DeleteJPNICTech == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("JPNICTech").Delete(c.JPNICTech[0])
	} else if service.DeleteConnection == base {
		err = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("IPv4").Delete(c.Connection[0])
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return err
}

func Get(base int, data *core.Service) service.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return service.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return service.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var serviceStruct []core.Service

	if base == service.ID { //ID
		err = db.Preload("IP").
			Preload("IP.Plan").
			Preload("Connection").
			Preload("Connection.NOC").
			Preload("Connection.BGPRouter").
			Preload("Connection.TunnelEndPointRouterIP").
			Preload("JPNICAdmin").
			Preload("JPNICTech").
			Preload("Group").
			First(&serviceStruct, data.ID).Error
	} else if base == service.Org { //Mail
		err = db.Preload("IP").
			Preload("Connection").
			Preload("JPNICAdmin").
			Preload("JPNICTech").
			Where("org = ?", data.Org).Find(&serviceStruct).Error
	} else if base == service.GID {
		err = db.Preload("IP").
			Preload("Connection").
			Preload("JPNICAdmin").
			Preload("JPNICTech").
			Where("group_id = ?", data.GroupID).Find(&serviceStruct).Error
	} else if base == service.GIDAndAddAllow {
		err = db.Where("group_id = ? AND add_allow = ?", data.GroupID, true).Find(&serviceStruct).Error
	} else if base == service.SearchNewNumber {
		err = db.Where("group_id = ?", data.GroupID).Find(&serviceStruct).Error
	} else if base == service.Open {
		err = db.Where("group_id = ? AND open = ?", data.GroupID, true).
			Preload("IP", "open = ?", true).
			Preload("Connection", "open = ?", true).
			Preload("Connection.NOC").
			Preload("Connection.BGPRouter").
			Preload("Connection.TunnelEndPointRouterIP").
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
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return service.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var services []core.Service
	err = db.Preload("IP").
		Preload("Connection").
		Preload("Connection.NOC").
		Preload("Connection.BGPRouter").
		Preload("Connection.TunnelEndPointRouterIP").
		Preload("JPNICAdmin").
		Preload("JPNICTech").
		Find(&services).Error
	return service.ResultDatabase{Err: err, Service: services}
}
