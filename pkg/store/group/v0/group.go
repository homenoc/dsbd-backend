package v0

import (
	"fmt"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(group *group.Group) (*group.Group, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return group, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&group).Error
	return group, err
}

func Delete(group *group.Group) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(group).Error
}

func Update(base int, g group.Group) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if group.UpdateOrg == base {
		result = db.Model(&group.Group{Model: gorm.Model{ID: g.ID}}).Update("org", g.Org)
	} else if group.UpdateStatus == base {
		result = db.Model(&group.Group{Model: gorm.Model{ID: g.ID}}).Update("status", g.Status)
	} else if group.UpdateInfo == base {
		result = db.Model(&group.Group{Model: gorm.Model{ID: g.ID}}).Update(group.Group{
			Name: g.Name, PostCode: g.PostCode, Address: g.Address, Mail: g.Mail, Phone: g.Phone, Country: g.Country})
	} else if group.UpdateAll == base {
		result = db.Model(&group.Group{Model: gorm.Model{ID: g.ID}}).Update(group.Group{
			Agree: g.Agree, Question: g.Question, Org: g.Org, Status: g.Status, Bandwidth: g.Bandwidth,
			Monitor: g.Monitor, Contract: g.Contract, Name: g.Name, PostCode: g.PostCode, Address: g.Address,
			Mail: g.Mail, Phone: g.Phone, Country: g.Country, Comment: g.Comment})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *group.Group) group.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var groupStruct []group.Group

	if base == group.ID { //ID
		err = db.First(&groupStruct, data.ID).Error
	} else if base == group.Org { //Org
		err = db.Where("org = ?", data.Org).Find(&groupStruct).Error
	} else if base == group.Email { //Mail
		err = db.Where("mail = ?", data.Mail).Find(&groupStruct).Error
	} else {
		log.Println("base select error")
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return group.ResultDatabase{Group: groupStruct, Err: err}
}

func GetAll() group.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return group.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var groups []group.Group
	err = db.Find(&groups).Error
	return group.ResultDatabase{Group: groups, Err: err}
}
