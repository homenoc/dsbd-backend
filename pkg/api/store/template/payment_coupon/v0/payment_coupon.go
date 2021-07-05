package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
	"log"
	"time"
)

func Create(coupon *core.PaymentCouponTemplate) (*core.PaymentCouponTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database coupon error")
		return coupon, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Create(&coupon).Error
	return coupon, err
}

func Delete(coupon *core.PaymentCouponTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database coupon error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	return db.Delete(coupon).Error
}

func Update(c core.PaymentCouponTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database coupon error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = nil

	err = db.Model(&core.PaymentCouponTemplate{Model: gorm.Model{ID: c.ID}}).Updates(core.PaymentCouponTemplate{
		Title:        c.Title,
		DiscountRate: c.DiscountRate,
		Comment:      c.Comment,
	}).Error

	return err
}

func Get(id uint) (core.PaymentCouponTemplate, error) {
	var coupon core.PaymentCouponTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database coupon error")
		return coupon, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return coupon, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.First(&coupon, id).Error

	return coupon, err
}

func GetAll() ([]core.PaymentCouponTemplate, error) {
	var coupons []core.PaymentCouponTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database coupon error")
		return coupons, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	dbSQL, err := db.DB()
	if err != nil {
		log.Printf("database error: %v", err)
		return coupons, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer dbSQL.Close()

	err = db.Find(&coupons).Error

	return coupons, err
}
