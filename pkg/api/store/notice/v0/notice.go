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

// Update writes the admin-editable columns of the notice row. Value-typed
// columns are always written (so clearing persists); pointer and time columns
// only when provided, so an omitted field never nulls the row.
func Update(data core.Notice) error {
	cols := []string{"title", "data"}
	if !data.StartTime.IsZero() {
		cols = append(cols, "start_time")
	}
	if !data.EndTime.IsZero() {
		cols = append(cols, "end_time")
	}
	if data.Everyone != nil {
		cols = append(cols, "everyone")
	}
	if data.Important != nil {
		cols = append(cols, "important")
	}
	if data.Fault != nil {
		cols = append(cols, "fault")
	}
	if data.Info != nil {
		cols = append(cols, "info")
	}
	return store.DB().Model(&core.Notice{Model: gorm.Model{ID: data.ID}}).Select(cols).Updates(data).Error
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
