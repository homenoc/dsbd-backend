package store

import (
	"github.com/homenoc/dsbd-backend/pkg/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
)

func ConnectDB() (*gorm.DB, error) {
	user := config.Conf.DB.User
	pass := config.Conf.DB.Pass
	protocol := "tcp(" + config.Conf.DB.IP + ":" + strconv.Itoa(config.Conf.DB.Port) + ")"
	dbName := config.Conf.DB.DBName

	db, err := gorm.Open("mysql", user+":"+pass+"@"+protocol+"/"+dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
