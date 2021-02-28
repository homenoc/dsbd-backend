package store

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicAdmin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnicTech"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/gateway"
	nocRouter "github.com/homenoc/dsbd-backend/pkg/api/core/noc/router"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/chat"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
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
	result := db.AutoMigrate(
		&user.User{},
		&group.Group{},
		&token.Token{},
		&network.Network{},
		&network.IP{},
		&network.JPNICAdmin{},
		&network.JPNICTech{},
		&connection.Connection{},
		&jpnicAdmin.JpnicAdmin{},
		&jpnicTech.JpnicTech{},
		&notice.Notice{},
		&ticket.Ticket{},
		&chat.Chat{},
		&noc.NOC{},
		&gateway.Gateway{},
		&nocRouter.Router{},
	)
	log.Println(result.Error)
}
