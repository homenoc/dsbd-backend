package store

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
)

func ConnectDB() (*gorm.DB, error) {
	user := config.Conf.DB.User
	pass := config.Conf.DB.Pass
	protocol := "tcp(" + config.Conf.DB.IP + ":" + strconv.Itoa(config.Conf.DB.Port) + ")"
	dbName := config.Conf.DB.DBName

	dsn := user + ":" + pass + "@" + protocol + "/" + dbName + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
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
		&core.MailTemplate{},
	)
	log.Println(result)
}
