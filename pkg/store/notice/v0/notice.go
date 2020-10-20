package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(notice *notice.Notice) (*notice.Notice, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return notice, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&notice).Error
	return notice, err
}

func Delete(notice *notice.Notice) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(notice).Error
}

func Update(base int, data notice.Notice) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if notice.UpdateAll == base {
		result = db.Model(&notice.Notice{Model: gorm.Model{ID: data.ID}}).Update(notice.Notice{
			UserID: data.UserID, GroupID: data.GroupID, StartTime: data.StartTime, EndingTime: data.EndingTime,
			Important: data.Important, Fault: data.Fault, Info: data.Info, Title: data.Title, Data: data.Data})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *notice.Notice) notice.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return notice.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var noticeStruct []notice.Notice

	if base == notice.ID { //ID
		err = db.First(&noticeStruct, data.ID).Error
	} else if base == notice.UserID { //UserID
		err = db.Where("user_id = ?", data.UserID).Find(&noticeStruct).Error
	} else if base == notice.GroupID { //GroupID
		err = db.Where("group_id = ?", data.GroupID).Find(&noticeStruct).Error
	} else if base == notice.UserIDAndGroupID { //UserID And GroupID
		err = db.Where("user_id = ? AND group_id = ?", data.UserID, data.GroupID).Find(&noticeStruct).Error
	} else if base == notice.Data { //Data
		err = db.Where("everyone = ?", data.Everyone).Or(notice.Notice{UserID: data.UserID, GroupID: 0}).
			Or(notice.Notice{GroupID: data.GroupID}).Order("id asc").Find(&noticeStruct).Error
	} else if base == notice.Everyone { //Everyone
		err = db.Where("everyone = ?", data.Everyone).Find(&noticeStruct).Error
	} else if base == notice.Important { //Important
		err = db.Where("important = ?", data.Important).Find(&noticeStruct).Error
	} else if base == notice.Fault { //Fault
		err = db.Where("fault = ?", data.Fault).Find(&noticeStruct).Error
	} else if base == notice.Info { //Info
		err = db.Where("info = ?", data.Info).Find(&noticeStruct).Error
	} else {
		log.Println("base select error")
		return notice.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return notice.ResultDatabase{Notice: noticeStruct, Err: err}
}

func GetAll() notice.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return notice.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var notices []notice.Notice
	err = db.Find(&notices).Error
	return notice.ResultDatabase{Notice: notices, Err: err}
}
