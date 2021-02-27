package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(u *user.User) error {
	result := Get(user.Email, &user.User{Email: u.Email})
	if result.Err != nil {
		return result.Err
	}
	if len(result.User) != 0 {
		log.Println("error: this email is already registered: " + u.Email)
		return fmt.Errorf("error: this email is already registered")
	}

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}

	defer db.Close()

	return db.Create(&u).Error
}

func Delete(u *user.User) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(u).Error
}

func Update(base int, u *user.User) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if user.UpdateVerifyMail == base {
		result = db.Model(&user.User{Model: gorm.Model{ID: u.ID}}).Update(user.User{MailVerify: u.MailVerify})
	} else if user.UpdateInfo == base {
		result = db.Model(&user.User{Model: gorm.Model{ID: u.ID}}).Update(user.User{
			Name:       u.Name,
			NameEn:     u.NameEn,
			Email:      u.Email,
			Pass:       u.Pass,
			Tech:       u.Tech,
			Level:      u.Level,
			MailVerify: u.MailVerify,
			MailToken:  u.MailToken,
			Org:        u.Org,
			OrgEn:      u.OrgEn,
			PostCode:   u.PostCode,
			Address:    u.Address,
			AddressEn:  u.AddressEn,
			Dept:       u.Dept,
			DeptEn:     u.DeptEn,
			Pos:        u.Pos,
			PosEn:      u.PosEn,
			Tel:        u.Tel,
			Fax:        u.Fax,
			Country:    u.Country,
		})
	} else if user.UpdateStatus == base {
		result = db.Model(&user.User{Model: gorm.Model{ID: u.ID}}).Update(user.User{Status: u.Status})
	} else if user.UpdateGID == base {
		result = db.Model(&user.User{Model: gorm.Model{ID: u.ID}}).Update(user.User{GroupID: u.GroupID})
	} else if user.UpdateLevel == base {
		result = db.Model(&user.User{Model: gorm.Model{ID: u.ID}}).Update("level", u.Level)
	} else if user.UpdateAll == base {
		result = db.Model(&user.User{Model: gorm.Model{ID: u.ID}}).Update(user.User{
			GroupID:    u.GroupID,
			Name:       u.Name,
			NameEn:     u.NameEn,
			Email:      u.Email,
			Pass:       u.Pass,
			Tech:       u.Tech,
			Level:      u.Level,
			MailVerify: u.MailVerify,
			MailToken:  u.MailToken,
			Org:        u.Org,
			OrgEn:      u.OrgEn,
			PostCode:   u.PostCode,
			Address:    u.Address,
			AddressEn:  u.AddressEn,
			Dept:       u.Dept,
			DeptEn:     u.DeptEn,
			Pos:        u.Pos,
			PosEn:      u.PosEn,
			Tel:        u.Tel,
			Fax:        u.Fax,
			Country:    u.Country,
			Status:     u.Status,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return result.Error
}

// value of base can reference from api/core/user/interface.go
func Get(base int, u *user.User) user.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var userStruct []user.User

	if base == user.ID { //ID
		err = db.First(&userStruct, u.ID).Error
	} else if base == user.GID { //GroupID
		err = db.Where("group_id = ?", u.GroupID).Find(&userStruct).Error
	} else if base == user.Email { //Mail
		err = db.Where("email = ?", u.Email).First(&userStruct).Error
	} else if base == user.MailToken { //Token
		err = db.Where("mail_token = ?", u.MailToken).Find(&userStruct).Error
	} else if base == user.GIDAndLevel { //GroupID and Level
		err = db.Where("group_id = ? AND level = ?", u.GroupID, u.Level).Find(&userStruct).Error
	} else {
		log.Println("base select error")
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}

	return user.ResultDatabase{User: userStruct, Err: err}
}

func GetAll() user.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var users []user.User
	err = db.Find(&users).Error
	return user.ResultDatabase{User: users, Err: err}
}
