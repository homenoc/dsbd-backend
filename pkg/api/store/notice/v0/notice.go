package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

//
// DBに入っている情報はUTCベースなので注意が必要
//

func Create(notice *core.Notice) (*core.Notice, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return notice, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&notice).Error
	return notice, err
}

func Delete(notice *core.Notice) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(notice).Error
}

func Update(base int, data core.Notice) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if notice.UpdateAll == base {
		result = db.Model(&core.Notice{Model: gorm.Model{ID: data.ID}}).Update(core.Notice{
			UserID:    data.UserID,
			GroupID:   data.GroupID,
			NOCID:     data.NOCID,
			StartTime: data.StartTime,
			EndTime:   data.EndTime,
			Important: data.Important,
			Everyone:  data.Everyone,
			Fault:     data.Fault,
			Info:      data.Info,
			Title:     data.Title,
			Data:      data.Data,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *core.Notice) notice.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return notice.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var noticeStruct []core.Notice

	//DBに入っているデータがUTCベースのため
	dateTime := time.Now().Add(9 * time.Hour)

	if base == notice.ID { //ID
		err = db.First(&noticeStruct, data.ID).Error
	} else if base == notice.UIDOrAll { //UserID Or All
		err = db.Where("user_id = ? AND start_time < ? AND ? < end_time", data.UserID, dateTime, dateTime).
			Or("everyone = ? AND start_time < ? AND ? < end_time", true, dateTime, dateTime).
			Order("id asc").Find(&noticeStruct).Error
	} else if base == notice.UIDOrGIDOrAll { //UserID Or GroupID Or All
		err = db.Where("user_id = ? AND start_time < ? AND ? < end_time", data.UserID, dateTime, dateTime).
			Or("group_id = ? AND start_time < ? AND ? < end_time", data.GroupID, dateTime, dateTime).
			Or("everyone = ? AND start_time < ? AND ? < end_time", true, dateTime, dateTime).
			Order("id asc").Find(&noticeStruct).Error
	} else if base == notice.UIDOrGIDOrNOCAllOrAll { //UserID Or GroupID Or NOCAll Or All
		err = db.Where("user_id = ? AND start_time < ? AND ? < end_time", data.UserID, dateTime, dateTime).
			Or("group_id = ? AND start_time < ? AND ? < end_time", data.GroupID, dateTime, dateTime).
			Or("everyone = ? AND start_time < ? AND ? < end_time", true, dateTime, dateTime).
			Or("noc_id != ? AND start_time < ? AND ? < end_time", 0, dateTime, dateTime).
			Order("id asc").Find(&noticeStruct).Error
	} else if base == notice.NOCAll { //UserID Or GroupID Or NOCAll Or All
		err = db.Where("user_id = ? AND user_id = ? AND noc_id != ? AND start_time < ? AND ? < end_time ",
			0, 0, 0, dateTime, dateTime).
			Order("id asc").Find(&noticeStruct).Error
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

func GetArray(base int, data *core.Notice, array []string) notice.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return notice.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var noticeStruct []core.Notice

	dateTime := time.Now()

	if base == notice.UIDOrGIDOrNOCAllOrAll { //UserID Or GroupID Or NOCAll Or All
		err = db.Where("user_id = ? AND start_time < ? AND ? < end_time ", data.UserID, dateTime, dateTime).
			Or("group_id = ? AND start_time < ? AND ? < end_time", data.GroupID, dateTime, dateTime).
			Or("everyone = ? AND start_time < ? AND ? < end_time", true, dateTime, dateTime).
			Or("noc_id IN (?) AND start_time < ? AND ? < end_time", array, dateTime, dateTime).
			Order("id desc").
			Find(&noticeStruct).Error
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

	var notices []core.Notice
	err = db.Find(&notices).Error
	return notice.ResultDatabase{Notice: notices, Err: err}
}
