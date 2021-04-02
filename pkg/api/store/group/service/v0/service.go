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

	var result error

	if service.UpdateData == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Update(core.Service{
			Org:       c.Org,
			OrgEn:     c.OrgEn,
			PostCode:  c.PostCode,
			Address:   c.Address,
			AddressEn: c.AddressEn,
			RouteV4:   c.RouteV4,
			RouteV6:   c.RouteV6,
			ASN:       c.ASN,
		}).Error
	} else if service.UpdateRoute == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Update(core.Service{
			RouteV4: c.RouteV4, RouteV6: c.RouteV6}).Error
	} else if service.UpdateGID == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Update(core.Service{GroupID: c.GroupID}).Error
	} else if service.UpdateStatus == base {
	} else if service.UpdateAll == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Update(c).Error
	} else if service.ReplaceIP == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("IP").Replace(c.IP[0]).Error
	} else if service.AppendIP == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("IP").Replace(c.IP[0]).Error
	} else if service.AppendJPNICAdmin == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("JPNICAdmin").Append(c.JPNICAdmin).Error
	} else if service.AppendJPNICTech == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("JPNICTech").Append(c.JPNICTech[0]).Error
	} else if service.AppendConnection == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("IPv4").Replace(c.Connection[0]).Error
	} else if service.DeleteIP == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("IP").Replace(c.IP[0]).Error
	} else if service.DeleteJPNICAdmin == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("JPNICAdmin").Delete(c.JPNICAdmin).Error
	} else if service.DeleteJPNICTech == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("JPNICTech").Delete(c.JPNICTech[0]).Error
	} else if service.DeleteConnection == base {
		result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Association("IPv4").Delete(c.Connection[0]).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	result = db.Model(&core.Service{Model: gorm.Model{ID: c.ID}}).Update(core.Service{AddAllow: c.AddAllow}).Error
	return result
}

func JoinJPNICAdmin(serviceID uint, input core.JPNICAdmin) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	db.Model(&core.Service{Model: gorm.Model{ID: serviceID}}).
		Association("JPNICAdmin").
		Append(input)

	return nil
}

func UpdateJPNICAdmin(serviceID uint, input core.JPNICAdmin) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	db.Model(&core.Service{Model: gorm.Model{ID: serviceID}}).
		Association("JPNICAdmin").
		Replace(&core.JPNICAdmin{Model: gorm.Model{ID: input.ID}}, &input)

	return nil
}

func DeleteJPNICAdmin(serviceID, jpnicAdminID uint) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	db.Model(&core.Service{Model: gorm.Model{ID: serviceID}}).
		Association("JPNICAdmin").
		Delete(&core.JPNICAdmin{Model: gorm.Model{ID: jpnicAdminID}})

	return nil
}

func JoinJPNICTech(serviceID uint, input core.JPNICTech) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	db.Model(&core.Service{Model: gorm.Model{ID: serviceID}}).
		Association("JPNICTech").
		Append(&input)

	return nil
}

func UpdateJPNICTech(serviceID uint, before, after core.JPNICTech) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	after.Model = before.Model
	db.Model(&core.Service{Model: gorm.Model{ID: serviceID}}).
		Association("JPNICTech").
		Replace(&before, &after)

	return nil
}

func DeleteJPNICTech(serviceID, jpnicTechID uint) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	db.Model(&core.Service{Model: gorm.Model{ID: serviceID}}).
		Association("JPNICTech").
		Delete(&core.JPNICTech{Model: gorm.Model{ID: jpnicTechID}})

	return nil
}

func JoinIP(serviceID uint, input core.IP) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	db.Model(&core.Service{Model: gorm.Model{ID: serviceID}}).
		Association("IP").
		Append(&input)

	return nil
}

func DeleteIP(serviceID, jpnicTechID uint) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	db.Model(&core.Service{Model: gorm.Model{ID: serviceID}}).
		Association("IP").
		Delete(&core.IP{Model: gorm.Model{ID: jpnicTechID}})

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
			Preload("Group").
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
	err = db.Preload("IP").
		Preload("Connection").
		Preload("Connection.ConnectionTemplate").
		Preload("Connection.NOC").
		Preload("Connection.BGPRouter").
		Preload("Connection.TunnelEndPointRouterIP").
		Preload("ServiceTemplate").
		Preload("JPNICAdmin").
		Preload("JPNICTech").
		Find(&services).Error
	return service.ResultDatabase{Err: err, Service: services}
}
