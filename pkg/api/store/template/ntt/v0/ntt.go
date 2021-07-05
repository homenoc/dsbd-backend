package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	ntt "github.com/homenoc/dsbd-backend/pkg/api/core/template/ntt"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"log"
	"time"
)

func Create(ntt *core.NTTTemplate) (*core.NTTTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return ntt, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&ntt).Error
	return ntt, err
}

func Delete(ntt *core.NTTTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	return db.Delete(ntt).Error
}

func Update(base int, c core.NTTTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = nil

	return err
}

func Get(base int, data *core.NTTTemplate) ntt.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return ntt.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ntt.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var nttStruct []core.NTTTemplate

	if base == ntt.ID { //ID
		err = db.First(&nttStruct, data.ID).Error
	} else {
		log.Println("base select error")
		return ntt.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return ntt.ResultDatabase{NTTs: nttStruct, Err: err}
}

func GetAll() ntt.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return ntt.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ntt.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var ntts []core.NTTTemplate
	err = db.Find(&ntts).Error
	return ntt.ResultDatabase{NTTs: ntts, Err: err}

}
