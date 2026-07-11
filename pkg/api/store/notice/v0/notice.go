package v0

import (
	"time"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

//
// DBに入っている情報はUTCベースなので注意が必要
//

func Create(n *core.Notice) (*core.Notice, error) {
	err := store.DB().Create(&n).Error
	return n, err
}

func Delete(n *core.Notice) error {
	return store.DB().Delete(n).Error
}

// UpdateAll updates the editable notice fields (was Update(UpdateAll, ...)).
func UpdateAll(data core.Notice) error {
	return store.DB().Model(&core.Notice{Model: gorm.Model{ID: data.ID}}).Updates(core.Notice{
		StartTime: data.StartTime,
		EndTime:   data.EndTime,
		Important: data.Important,
		Everyone:  data.Everyone,
		Fault:     data.Fault,
		Info:      data.Info,
		Title:     data.Title,
		Data:      data.Data,
	}).Error
}

// GetByID looks up one notice by primary key.
func GetByID(id uint) notice.ResultDatabase {
	var notices []core.Notice
	err := store.DB().First(&notices, id).Error
	return notice.ResultDatabase{Notice: notices, Err: err}
}

// GetActiveForUser returns notices currently in their display window that
// target the given user or everyone (was Get(UIDOrAll, ...)). DB times are UTC.
func GetActiveForUser(userID uint) notice.ResultDatabase {
	dateTime := time.Now().Add(9 * time.Hour)
	var notices []core.Notice
	err := store.DB().Where("start_time < ? AND ? < end_time", dateTime, dateTime).
		Joins("left outer join notice_user on notices.id = notice_user.notice_id").
		Where("notice_user.user_id = ?", userID).
		Or("everyone = ? AND start_time < ? AND ? < end_time", true, dateTime, dateTime).
		Order("id asc").Find(&notices).Error
	return notice.ResultDatabase{Notice: notices, Err: err}
}

func GetAll() notice.ResultDatabase {
	var notices []core.Notice
	err := store.DB().Find(&notices).Error
	return notice.ResultDatabase{Notice: notices, Err: err}
}
