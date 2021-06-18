package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func Create(coupon *core.PaymentCouponTemplate) (*core.PaymentCouponTemplate, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database coupon error")
		return coupon, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&coupon).Error
	return coupon, err
}

func Delete(coupon *core.PaymentCouponTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database coupon error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(coupon).Error
}

func Update(c core.PaymentCouponTemplate) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database coupon error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	result = db.Model(&core.PaymentCouponTemplate{Model: gorm.Model{ID: c.ID}}).Update(core.PaymentCouponTemplate{
		Title:        c.Title,
		DiscountRate: c.DiscountRate,
		Comment:      c.Comment,
	})

	return result.Error
}

func Get(id uint) (core.PaymentCouponTemplate, error) {
	var coupon core.PaymentCouponTemplate

	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database coupon error")
		return coupon, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

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
	defer db.Close()

	err = db.Find(&coupons).Error

	return coupons, err
}
