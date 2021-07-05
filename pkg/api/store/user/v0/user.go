package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(u *core.User) error {
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

	return db.Create(&u).Error
}

func Delete(u *core.User) error {
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

	return db.Delete(u).Error
}

func Update(base int, u *core.User) error {
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

	if user.UpdateVerifyMail == base {
		err = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Updates(core.User{MailVerify: u.MailVerify}).Error
	} else if user.UpdateInfo == base {
		err = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Updates(core.User{
			Name:       u.Name,
			NameEn:     u.NameEn,
			Email:      u.Email,
			Pass:       u.Pass,
			Level:      u.Level,
			MailVerify: u.MailVerify,
			MailToken:  u.MailToken,
		}).Error
	} else if user.UpdateGID == base {
		err = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Updates(core.User{GroupID: u.GroupID}).Error
	} else if user.UpdateLevel == base {
		err = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Updates(core.User{Level: u.Level}).Error
	} else if user.UpdateAll == base {
		err = db.Model(&core.User{Model: gorm.Model{ID: u.ID}}).Updates(core.User{
			GroupID:       u.GroupID,
			Name:          u.Name,
			NameEn:        u.NameEn,
			Email:         u.Email,
			Pass:          u.Pass,
			Level:         u.Level,
			MailVerify:    u.MailVerify,
			MailToken:     u.MailToken,
			ExpiredStatus: u.ExpiredStatus,
		}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	return err
}

// value of base can reference from api/core/user/interface.go
func Get(base int, u *core.User) user.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var userStruct []core.User

	if base == user.ID { //ID
		err = db.First(&userStruct, u.ID).Error
	} else if base == user.IDDetail {
		err = db.Where("id = ?", u.ID).
			Preload("Ticket").
			Preload("Ticket.Chat").
			Preload("Group").
			Preload("Group.PaymentMembershipTemplate").
			Preload("Group.PaymentCouponTemplate").
			Preload("Group.Users").
			Preload("Group.Services").
			Preload("Group.Tickets").
			Preload("Group.Tickets.Chat").
			Preload("Group.Services.IP").
			Preload("Group.Services.IP.Plan").
			Preload("Group.Services.Connection").
			Preload("Group.Services.Connection.ConnectionTemplate").
			Preload("Group.Services.Connection.NOC").
			Preload("Group.Services.Connection.BGPRouter").
			Preload("Group.Services.Connection.BGPRouter.NOC").
			Preload("Group.Services.Connection.TunnelEndPointRouterIP").
			Preload("Group.Services.ServiceTemplate").
			Preload("Group.Services.JPNICAdmin").
			Preload("Group.Services.JPNICTech").Find(&userStruct).Error
	} else if base == user.GID { //GroupID
		err = db.Where("group_id = ?", u.GroupID).Find(&userStruct).Error
	} else if base == user.Email { //Mail
		err = db.Where("email = ?", u.Email).First(&userStruct).Error
	} else if base == user.MailToken { //Token
		err = db.Where("mail_token = ?", u.MailToken).Find(&userStruct).Error
	} else if base == user.GIDAndLevel { //GroupID and Level
		err = db.Where("group_id = ? AND level = ?", u.GroupID, u.Level).Find(&userStruct).Error
	} else if base == user.IDGetGroup { //GroupID and Level
		err = db.First(&userStruct, u.ID).
			Preload("Group").
			Find(&userStruct).Error

	} else {
		log.Println("base select error")
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}

	return user.ResultDatabase{User: userStruct, Err: err}
}

// value of base can reference from api/core/user/interface.go
func GetArray(u []uint) user.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var userStruct []core.User

	err = db.Where(u).Find(&userStruct).Error

	return user.ResultDatabase{User: userStruct, Err: err}
}

func GetAll() user.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}

	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return user.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var users []core.User
	err = db.Find(&users).Error
	return user.ResultDatabase{User: users, Err: err}
}
