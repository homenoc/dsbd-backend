package store

import (
	"log"
	"strconv"
	"time"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// db is the process-wide connection pool, opened once by Init. Store functions
// obtain it via DB(). Previously every store call opened a fresh gorm.Open and
// closed the underlying *sql.DB with defer, which destroyed pooling on every
// query; the pool now lives for the process lifetime.
var db *gorm.DB

// testDB, when set, overrides db (for tests). Kept for the existing SetTestDB seam.
var testDB *gorm.DB

// SetTestDB sets a test database connection (for testing only).
func SetTestDB(d *gorm.DB) {
	testDB = d
}

// ClearTestDB clears the test database connection.
func ClearTestDB() {
	testDB = nil
}

// Init opens the shared connection pool. Safe to call more than once; a test DB
// injected via SetTestDB always wins in DB().
func Init() error {
	if db != nil {
		return nil
	}
	opened, err := open()
	if err != nil {
		return err
	}
	db = opened
	return nil
}

func open() (*gorm.DB, error) {
	user := config.Conf.DB.User
	pass := config.Conf.DB.Pass
	protocol := "tcp(" + config.Conf.DB.IP + ":" + strconv.Itoa(config.Conf.DB.Port) + ")"
	dbName := config.Conf.DB.DBName

	dsn := user + ":" + pass + "@" + protocol + "/" + dbName + "?charset=utf8&parseTime=True&loc=Local"
	opened, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := opened.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return opened, nil
}

// DB returns the shared pool (or the injected test DB). It lazily initializes
// the pool on first use so callers that never called Init still work.
func DB() *gorm.DB {
	if testDB != nil {
		return testDB
	}
	if db == nil {
		if err := Init(); err != nil {
			log.Fatalf("[DB] Failed to connect: %v", err)
		}
	}
	return db
}

// Tx runs fn inside a transaction on the shared pool. Store write functions
// accept an optional *gorm.DB so they can be composed inside a Tx.
func Tx(fn func(tx *gorm.DB) error) error {
	return DB().Transaction(fn)
}

// Migrate connects (via Init) and runs AutoMigrate for every model. It is the
// explicit schema-setup step used by the `init database` command and tests —
// deliberately separate from Init(), which servers call at boot so that
// starting a server does not re-run migrations.
func Migrate() {
	log.Println("[DB] Connecting to database...")
	if err := Init(); err != nil {
		log.Fatalf("[DB] Failed to connect: %v", err)
	}
	log.Println("[DB] Connected successfully")

	log.Println("[DB] Running migrations...")
	err := DB().AutoMigrate(
		&core.User{},
		&core.Group{},
		&core.Memo{},
		&core.Service{},
		&core.Connection{},
		&core.NOC{},
		&core.BGPRouter{},
		&core.TunnelEndPointRouter{},
		&core.TunnelEndPointRouterIP{},
		&core.IP{},
		&core.Plan{},
		&core.JPNICAdmin{},
		&core.JPNICTech{},
		&core.Ticket{},
		&core.Chat{},
		&core.Token{},
		&core.Notice{},
	)
	if err != nil {
		log.Fatalf("[DB] Migration failed: %v", err)
	}
	log.Println("[DB] Migration completed successfully")
}
