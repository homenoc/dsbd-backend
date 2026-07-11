package v0

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	store.SetTestDB(gormDB)

	cleanup := func() {
		store.ClearTestDB()
		db.Close()
	}

	return mock, cleanup
}

func TestUpdatePartField(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `nocs`").
		WithArgs(sqlmock.AnyArg(), "nocTest", uint(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	testTemplate := core.NOC{
		Model: gorm.Model{ID: 1},
		Name:  "nocTest",
	}

	if err := UpdateAll(testTemplate); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
