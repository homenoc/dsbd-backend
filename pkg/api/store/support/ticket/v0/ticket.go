package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(support *ticket.Ticket) (*ticket.Ticket, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return support, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&support).Error
	return support, err
}

func Delete(support *ticket.Ticket) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(support).Error
}

func Update(base int, t ticket.Ticket) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	//#4 Issue(解決済み）
	if ticket.UpdateAll == base {
		result = db.Model(&ticket.Ticket{Model: gorm.Model{ID: t.ID}}).Update(&ticket.Ticket{Title: t.Title,
			GroupID: t.GroupID, UserID: t.UserID, ChatIDStart: t.ChatIDStart, ChatIDEnd: t.ChatIDEnd, Solved: t.Solved})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *ticket.Ticket) ticket.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return ticket.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var ticketStruct []ticket.Ticket

	if base == ticket.ID { //ID
		err = db.First(&ticketStruct, data.ID).Error
	} else if base == ticket.GID { //GroupID
		err = db.Where("group_id = ?", data.GroupID).Find(&ticketStruct).Error
	} else if base == ticket.UID { //UserID
		err = db.Where("user_id = ?", data.UserID).Find(&ticketStruct).Error
	} else if base == ticket.CIDStart { //ChatID Start
		err = db.Where("ticket_id = ?", data.ChatIDStart).Find(&ticketStruct).Error
	} else if base == ticket.CIDEnd { //ChatID End
		err = db.Where("ticket_id = ?", data.ChatIDEnd).Find(&ticketStruct).Error
	} else {
		log.Println("base select error")
		return ticket.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return ticket.ResultDatabase{Ticket: ticketStruct, Err: err}
}

func GetAll() ticket.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return ticket.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var tickets []ticket.Ticket
	err = db.Find(&tickets).Error
	return ticket.ResultDatabase{Ticket: tickets, Err: err}
}
