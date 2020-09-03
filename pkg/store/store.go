package store

import (
	"database/sql"
	"fmt"
	"git.bgp.ne.jp/dsbd/backend/pkg/config"
	"log"
	"strconv"
)

func connectDB() *sql.DB {
	conf, err := config.GetConfig()
	if err != nil {
		return nil
	}
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s", conf.DB.User, conf.DB.Pass, conf.DB.IP+strconv.Itoa(conf.DB.Port), conf.DB.DBName))
	if err != nil {
		log.Println("database ping error: ", err)
		return nil
	}
	if err := db.Ping(); err != nil {
		log.Println("database ping error: ", err)
		return nil
	}
	return db
}

func createDB(database string) error {
	db := *connectDB()
	_, err := db.Exec(database)
	if err != nil {
		log.Println("database create error: ", err)
		return err
	}
	return nil
}

func InitDB() error {
	// User data
	err := createDB(`CREATE TABLE IF NOT EXISTS "user" ("id" INTEGER PRIMARY KEY, "gid" INT, "name" VARCHAR(255), "level" INT, "email" VARCHAR(2000), "pass" VARCHAR(2000))`)
	if err != nil {
		log.Println("create error: User database ", err)
		return err
	}
	// Group data
	createDB(`CREATE TABLE IF NOT EXISTS "group" ("id" INTEGER PRIMARY KEY, "orgJa" VARCHAR(2000), "org" VARCHAR(2000), "level" INT)`)
	if err != nil {
		log.Println("create error: Group database ", err)
		return err
	}
	// Inquiry data
	createDB(`CREATE TABLE IF NOT EXISTS "inquiry" ("id" INTEGER PRIMARY KEY, "name" VARCHAR(255), "pass" VARCHAR(255))`)
	if err != nil {
		log.Println("create error: Inquiry database ", err)
		return err
	}
	// IPAssign data
	createDB(`CREATE TABLE IF NOT EXISTS "ip_assign" ("id" INTEGER PRIMARY KEY, "name" VARCHAR(255),"admin" VARCHAR(500),"user" VARCHAR(2000),"uuid" VARCHAR(20000),"maxvm" INT,"maxcpu" INT,"maxmem" INT,"maxstorage" INT,"net" VARCHAR(255))`)
	if err != nil {
		log.Println("create error: IPAssign database ", err)
		return err
	}
	// Token data
	createDB(`CREATE TABLE IF NOT EXISTS "token" ("id" INTEGER PRIMARY KEY, "uid" INT, "token" VARCHAR(255),"begintime" INT,"endtime" INT)`)
	if err != nil {
		log.Println("create error: Token database ", err)
		return err
	}
	// Administrator data
	createDB(`CREATE TABLE IF NOT EXISTS "administrator" ("id" INTEGER PRIMARY KEY, "uid" VARCHAR(255), "name" VARCHAR(255), "email" INT, "pass" INT)`)
	if err != nil {
		log.Println("create error: Administrator database ", err)
		return err
	}
	return nil
}
