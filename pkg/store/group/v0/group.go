package v0

import (
	"fmt"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(group *group.Group) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Create(group).Error
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
	} else if group.UpdateTechID == base {
		result = db.Model(&group.Group{Model: gorm.Model{ID: g.ID}}).Update("tech_id", g.TechID)
	} else if group.UpdateInfo == base {
		result = db.Model(&group.Group{Model: gorm.Model{ID: g.ID}}).Update(group.Group{
			Name: g.Name, PostCode: g.PostCode, Address: g.Address, Mail: g.Mail, Phone: g.Phone, Country: g.Country})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *group.Group) *gorm.DB {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return &gorm.DB{Error: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var groupStruct group.Group

	if base == group.ID { //ID
		return db.First(&groupStruct, group.ID)
	} else if base == group.Org { //Org
		return db.Where("org = ?", group.Org).Find(&groupStruct)
	} else if base == group.Email { //Mail
		return db.Where("email = ?", group.Email).Find(&groupStruct)
	} else {
		log.Println("base select error")
		return &gorm.DB{Error: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
}

func GetAll() *gorm.DB {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return &gorm.DB{Error: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var users []group.Group
	return db.Find(&users)
}
