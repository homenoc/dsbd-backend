package store

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
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
		&core.User{},
		&core.Group{},
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
		&core.ServiceTemplate{},
		&core.ConnectionTemplate{},
		&core.NTTTemplate{},
		&core.Ticket{},
		&core.Chat{},
		&core.Token{},
		&core.Notice{},
		&core.IPv4Template{},
		&core.IPv6Template{},
		&core.IPv4RouteTemplate{},
		&core.IPv6RouteTemplate{},
		&core.Payment{},
		&core.PaymentCouponTemplate{},
		&core.PaymentMembershipTemplate{},
		&core.PaymentDonateTemplate{},
	)
	log.Println(result.Error)
}
