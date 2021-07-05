package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(support *core.Ticket) (*core.Ticket, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return support, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&support).Error
	return support, err
}

func Delete(support *core.Ticket) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	return db.Delete(support).Error
}

func Update(base int, t core.Ticket) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = nil

	//#4 Issue(解決済み）
	if ticket.UpdateAll == base {
		err = db.Model(&core.Ticket{Model: gorm.Model{ID: t.ID}}).Updates(&core.Ticket{Title: t.Title,
			GroupID:       t.GroupID,
			UserID:        t.UserID,
			Solved:        t.Solved,
			Request:       t.Request,
			RequestReject: t.RequestReject,
		}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return err
}

func Get(base int, data *core.Ticket) ticket.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return ticket.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ticket.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var ticketStruct []core.Ticket

	if base == ticket.ID { //ID
		err = db.Preload("User").
			Preload("Group").
			Preload("Chat").
			Preload("Chat.User").
			First(&ticketStruct, data.ID).Error
	} else if base == ticket.GID { //GroupID
		err = db.Where("group_id = ?", data.GroupID).
			Preload("User").
			Preload("Group").
			Preload("Chat").
			Preload("Chat.User").
			Find(&ticketStruct).Error
	} else if base == ticket.UID { //UserID
		err = db.Where("user_id = ?", data.UserID).Find(&ticketStruct).Error
	} else {
		log.Println("base select error")
		return ticket.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return ticket.ResultDatabase{Tickets: ticketStruct, Err: err}
}

func GetAll() ticket.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return ticket.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return ticket.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer dbSQL.Close()

	var tickets []core.Ticket
	err = db.Preload("User").
		Preload("Group").
		Preload("Chat").
		Preload("Chat.User").
		Find(&tickets).Error
	return ticket.ResultDatabase{Tickets: tickets, Err: err}
}
