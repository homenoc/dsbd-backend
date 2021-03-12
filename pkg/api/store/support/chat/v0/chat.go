package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/chat"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(support *core.Chat) (*core.Chat, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return support, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&support).Error
	return support, err
}

func Delete(support *core.Chat) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(support).Error
}

func Update(base int, s core.Chat) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	if chat.UpdateUserID == base {
		result = db.Model(&core.Chat{Model: gorm.Model{ID: s.ID}}).Update(core.Chat{UserID: s.UserID})
	} else if chat.UpdateAll == base {
		result = db.Model(&core.Chat{Model: gorm.Model{ID: s.ID}}).Update(core.Chat{
			TicketID: s.TicketID,
			UserID:   s.UserID,
			Admin:    s.Admin,
			Data:     s.Data,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *core.Chat) chat.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return chat.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var chatStruct []core.Chat

	if base == chat.ID { //ID
		err = db.First(&chatStruct, data.ID).Error
	} else if base == chat.TicketID { //TicketID
		err = db.Where("ticket_id = ?", data.TicketID).Order("id asc").Find(&chatStruct).Error
	} else {
		log.Println("base select error")
		return chat.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return chat.ResultDatabase{Chat: chatStruct, Err: err}
}

func GetAll() chat.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return chat.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var chats []core.Chat
	err = db.Find(&chats).Error
	return chat.ResultDatabase{Chat: chats, Err: err}
}
