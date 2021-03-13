package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	ntt "github.com/homenoc/dsbd-backend/pkg/api/core/template/ntt"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(ntt *core.NTTTemplate) (*core.NTTTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return ntt, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&ntt).Error
	return ntt, err
}

func Delete(ntt *core.NTTTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(ntt).Error
}

func Update(base int, c core.NTTTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	//if base == ntt.UpdateAll {
	//	result = db.Model(&core.NTTTemplate{Model: gorm.Model{ID: c.ID}}).Update(core.NTTTemplate{
	//		Hidden:  c.Hidden,
	//		Name:    c.Name,
	//		Comment: c.Comment,
	//	})
	//} else {
	//	log.Println("base select error")
	//	return fmt.Errorf("(%s)error: base select\n", time.Now())
	//}

	return result.Error
}

func Get(base int, data *core.NTTTemplate) ntt.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database ntt error")
		return ntt.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

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
	defer db.Close()

	var ntts []core.NTTTemplate
	err = db.Find(&ntts).Error
	return ntt.ResultDatabase{NTTs: ntts, Err: err}

}
