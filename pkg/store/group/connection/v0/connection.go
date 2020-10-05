package v0

import (
	"fmt"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(connection *connection.Connection) (*connection.Connection, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&connection).Error
	return connection, err
}

func Delete(connection *connection.Connection) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(connection).Error
}

func Update(base int, c connection.Connection) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if connection.UpdateInfo == base {
		result = db.Model(&connection.Connection{Model: gorm.Model{ID: c.ID}}).Update(connection.Connection{
			Service: c.Service, Fee: c.Fee, NOC: c.NOC, TermIP: c.TermIP, LinkV4: c.LinkV4, LinkV6: c.LinkV6})
	} else if connection.UpdateUserInfo == base {
		result = db.Model(&connection.Connection{Model: gorm.Model{ID: c.ID}}).Update(connection.Connection{
			NTT: c.NTT, Name: c.Name, Org: c.Org, PostCode: c.PostCode, Address: c.Address, Mail: c.Mail,
			Phone: c.Phone, Country: c.Country, Comment: c.Comment})
	} else if connection.UpdateGID == base {
		result = db.Model(&connection.Connection{Model: gorm.Model{ID: c.ID}}).Update(connection.Connection{GroupID: c.GroupID})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *connection.Connection) connection.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var connectionStruct []connection.Connection

	if base == connection.ID { //ID
		err = db.First(&connectionStruct, data.ID).Error
	} else if base == connection.Org { //Org
		err = db.Where("org = ?", data.Org).Find(&connectionStruct).Error
	} else if base == connection.Email { //Mail
		err = db.Where("mail = ?", data.Mail).Find(&connectionStruct).Error
	} else if base == connection.GID {
		err = db.Where("group_id = ?", data.GroupID).Find(&connectionStruct).Error
	} else {
		log.Println("base select error")
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return connection.ResultDatabase{Connection: connectionStruct, Err: err}
}

func GetAll() connection.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return connection.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var connections []connection.Connection
	err = db.Find(&connections).Error
	return connection.ResultDatabase{Connection: connections, Err: err}

}
