package store

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnic_user"
	networkUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/network_user"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/homenoc/dsbd-backend/pkg/tool/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"strconv"
)

func ConnectDB() (*gorm.DB, error) {
	user := config.Conf.DB.User
	pass := config.Conf.DB.Pass
	protocol := "tcp(" + config.Conf.DB.IP + ":" + strconv.Itoa(config.Conf.DB.Port) + ")"
	dbName := config.Conf.DB.DBName

	db, err := gorm.Open("mysql", user+":"+pass+"@"+protocol+"/"+dbName+"?parseTime=true")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB() {
	db, _ := ConnectDB()
	result := db.AutoMigrate(&user.User{}, &group.Group{}, &token.Token{}, &network.Network{}, &jpnic_user.JPNICUser{},
		&connection.Connection{}, &networkUser.NetworkUser{})
	log.Println(result.Error)
	//return nil
}
