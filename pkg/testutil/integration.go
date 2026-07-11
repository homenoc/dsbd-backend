// Package testutil provides helpers for integration tests that run the API
// against a real MySQL instance (make db-up).
package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/seed"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/slack-go/slack"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDBName is a dedicated schema so integration tests never touch the
// development database (dsbd-backend).
const TestDBName = "dsbd-backend-test"

// RepoRoot returns the repository root directory (resolved from this file).
func RepoRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(filepath.Dir(file)))
}

// SkipUnlessIntegration skips the test unless DSBD_TEST_DB=1 is set.
func SkipUnlessIntegration(tb testing.TB) {
	tb.Helper()
	if os.Getenv("DSBD_TEST_DB") != "1" {
		tb.Skip("integration test: set DSBD_TEST_DB=1 and run `make db-up` first")
	}
}

// SetupIntegration loads the repo config, points it at a freshly re-created
// test database, neutralizes external side effects (Slack/mail/log file),
// runs migrations and seeds deterministic data.
//
// Determinism note: the schema is dropped and re-created on every call, so
// auto-increment IDs produced by the seed are stable across runs.
func SetupIntegration(tb testing.TB) {
	tb.Helper()
	SkipUnlessIntegration(tb)

	if err := config.GetConfig(filepath.Join(RepoRoot(), "configs", "config.json")); err != nil {
		tb.Fatalf("failed to load config: %v", err)
	}

	// Isolate from the dev database and from external services.
	config.Conf.DB.DBName = TestDBName
	config.Conf.Mail.Host = "127.0.0.1"
	config.Conf.Mail.Port = 1 // unreachable: mail sends fail fast instead of dialing out
	config.Conf.Log.Path = filepath.Join(tb.TempDir(), "test.log")

	// Slack client pointed at an unroutable local endpoint: notification
	// calls fail fast without reaching the network.
	notification.Notification.Slack = slack.New("test-token",
		slack.OptionAPIURL("http://127.0.0.1:1/api/"))

	recreateTestDB(tb)

	// Inject a dedicated connection to the freshly-created test schema via the
	// SetTestDB seam so store.DB() uses it instead of the process pool.
	store.SetTestDB(openTestDB(tb))
	tb.Cleanup(func() { store.ClearTestDB() })

	store.Migrate()
	if err := seed.Run(); err != nil {
		tb.Fatalf("failed to seed test database: %v", err)
	}
}

func openTestDB(tb testing.TB) *gorm.DB {
	tb.Helper()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Conf.DB.User, config.Conf.DB.Pass, config.Conf.DB.IP,
		strconv.Itoa(config.Conf.DB.Port), TestDBName)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Discard,
	})
	if err != nil {
		tb.Fatalf("failed to open test DB connection: %v", err)
	}
	return conn
}

func recreateTestDB(tb testing.TB) {
	tb.Helper()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql?charset=utf8&parseTime=True&loc=Local",
		config.Conf.DB.User, config.Conf.DB.Pass, config.Conf.DB.IP, strconv.Itoa(config.Conf.DB.Port))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		tb.Fatalf("failed to connect to MySQL (is `make db-up` running?): %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		tb.Fatalf("failed to get sql.DB: %v", err)
	}
	defer sqlDB.Close()

	if err := db.Exec("DROP DATABASE IF EXISTS `" + TestDBName + "`").Error; err != nil {
		tb.Fatalf("failed to drop test database: %v", err)
	}
	if err := db.Exec("CREATE DATABASE `" + TestDBName + "`").Error; err != nil {
		tb.Fatalf("failed to create test database: %v", err)
	}
}
