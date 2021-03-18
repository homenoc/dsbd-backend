package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(u *core.User) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}

	defer db.Close()

	return db.Create(&u).Error
}

func Delete(u *core.User) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(u).Error
}

func Update(base int, u *core.User) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if user.UpdateVerifyMail == base {
		result = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Update(core.User{MailVerify: u.MailVerify})
	} else if user.UpdateInfo == base {
		result = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Update(core.User{
			Name:       u.Name,
			NameEn:     u.NameEn,
			Email:      u.Email,
			Pass:       u.Pass,
			Level:      u.Level,
			MailVerify: u.MailVerify,
			MailToken:  u.MailToken,
		})
	} else if user.UpdateGID == base {
		result = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Update(core.User{GroupID: u.GroupID})
	} else if user.UpdateLevel == base {
		result = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Update("level", u.Level)
	} else if user.UpdateAll == base {
		result = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Update(core.User{
			GroupID:       u.GroupID,
			Name:          u.Name,
			NameEn:        u.NameEn,
			Email:         u.Email,
			Pass:          u.Pass,
			Level:         u.Level,
			MailVerify:    u.MailVerify,
			MailToken:     u.MailToken,
			ExpiredStatus: u.ExpiredStatus,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return result.Error
}

// value of base can reference from api/core/user/interface.go
func Get(base int, u *core.User) user.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var userStruct []core.User

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

	var users []core.User
	err = db.Find(&users).Error
	return user.ResultDatabase{User: users, Err: err}
}
