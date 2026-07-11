package v0

import (
	"log"
	"time"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
	dbToken "github.com/homenoc/dsbd-backend/pkg/api/store/token/v0"
	"gorm.io/gorm"
)

func TokenRemove() {
	go func() {
		// 15分おき
		t := time.NewTicker(15 * 60 * time.Second)
		for {
			select {
			case <-t.C:
				expired, err := dbToken.GetExpired()
				if err != nil {
					log.Println(err)
				}
				for _, tmp := range expired {
					err := dbToken.Delete(&core.Token{Model: gorm.Model{ID: tmp.ID}})
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
		t.Stop() //停止
	}()
}
