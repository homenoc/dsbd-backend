package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbToken "github.com/homenoc/dsbd-backend/pkg/api/store/token/v0"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func TokenRemove() {
	go func() {
		// 15分おき
		t := time.NewTicker(15 * 60 * time.Second)
		for {
			select {
			case <-t.C:
				result := dbToken.Get(token.ExpiredTime, &core.Token{})
				if result.Err != nil {
					log.Println(result.Err)
				}
				for _, tmp := range result.Token {
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
