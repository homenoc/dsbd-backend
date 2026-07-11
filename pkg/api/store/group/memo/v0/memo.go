package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

func Create(connection *core.Memo) (*core.Memo, error) {
	db := store.DB()

	err := db.Create(&connection).Error
	return connection, err
}

func Delete(connection *core.Memo) error {
	db := store.DB()

	return db.Delete(connection).Error
}

func Update(memo core.Memo) error {
	db := store.DB()

	return db.Model(&core.Memo{Model: gorm.Model{ID: memo.ID}}).Updates(memo).Error
}
