package store

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/homenoc/dsbd-backend/pkg/config"
	"log"
	"strconv"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Conf.DB.User, config.Conf.DB.Pass, config.Conf.DB.IP+strconv.Itoa(config.Conf.DB.Port), config.Conf.DB.DBName))
	if err != nil {
		log.Println("database open error: ", err)
		return nil
	}
	if err := db.Ping(); err != nil {
		log.Println("database ping error: ", err)
		return nil
	}
	return db
}

func createDB(database string) error {
	db := *ConnectDB()
	_, err := db.Exec(database)
	if err != nil {
		log.Println("database create error: ", err)
		return err
	}
	return nil
}

func InitDB() error {
	// User data
	err := createDB(`CREATE TABLE IF NOT EXISTS "user" ("id" INTEGER PRIMARY KEY, "created_at" INT, "updated_at" INT,
"gid" INT,　"name" VARCHAR(255), "email" VARCHAR(2000), "pass" VARCHAR(2000), "level" INT, "is_verify" INT, "mail_token" VARCHAR(2000)`)
	if err != nil {
		log.Println("create error: User database ", err)
		return err
	}
	// Group data
	err = createDB(`CREATE TABLE IF NOT EXISTS "group" ("id" INTEGER PRIMARY KEY, "created_at" INT, "updated_at" INT,
"org_ja" VARCHAR(2000), "org" VARCHAR(2000), "status" INT, "tech_id" VARCHAR(2000), "postcode" VARCHAR(100), "address_ja" VARCHAR(2000),
"address" VARCHAR(2000), "mail" VARCHAR(2000), "phone" VARCHAR(2000))`)
	if err != nil {
		log.Println("create error: Group database ", err)
		return err
	}
	// private data
	err = createDB(`CREATE TABLE IF NOT EXISTS "private" ("id" INTEGER PRIMARY KEY, "name" VARCHAR(255), "pass" VARCHAR(255))`)
	if err != nil {
		log.Println("create error: Inquiry database ", err)
		return err
	}
	// IPAssign data
	err = createDB(`CREATE TABLE IF NOT EXISTS "ip_assign" ("id" INTEGER PRIMARY KEY, "name" VARCHAR(255), "admin" VARCHAR(500),"user" VARCHAR(2000),"uuid" VARCHAR(20000),"maxvm" INT,"maxcpu" INT,"maxmem" INT,"maxstorage" INT,"net" VARCHAR(255))`)
	if err != nil {
		log.Println("create error: IPAssign database ", err)
		return err
	}
	// Token data
	err = createDB(`CREATE TABLE IF NOT EXISTS "token" ("id" INTEGER PRIMARY KEY, "created_at" INT,
"updated_at" INT, "expired_at" INT, "deleted_at" INT, "uid" INT, "status" INT, "user_token" VARCHAR(1000),
"tmp_token" VARCHAR(1000), "access_token" VARCHAR(1000), "debug" VARCHAR(500))`)
	if err != nil {
		log.Println("create error: Token database ", err)
		return err
	}
	// Administrator data
	err = createDB(`CREATE TABLE IF NOT EXISTS "admin" ("id" INTEGER PRIMARY KEY, "name" VARCHAR(255), "email" INT, "pass" INT)`)
	if err != nil {
		log.Println("create error: Administrator database ", err)
		return err
	}
	return nil
}
