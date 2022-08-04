package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddByAdmin(c *gin.Context) {
	var input core.Group

	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if _, err = dbGroup.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.ResultAdmin{})
}

func DeleteByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbGroup.Delete(&core.Group{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.ResultAdmin{})
}

func UpdateByAdmin(c *gin.Context) {
	// ID取得
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// serviceIDが0の時エラー処理
	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("This id is wrong... ")})
		return
	}

	var input core.Group

	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	tmp := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: uint(id)}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: tmp.Err.Error()})
		return
	}

	noticeByAdmin(tmp.Group[0], input)

	input.ID = uint(id)
	if input.PaymentCouponTemplateID != nil && *input.PaymentCouponTemplateID == 0 {
		input.PaymentCouponTemplateID = nil
	}
	if input.PaymentMembershipTemplateID != nil && *input.PaymentMembershipTemplateID == 0 {
		input.PaymentMembershipTemplateID = nil
	}

	if err = dbGroup.Update(group.UpdateAll, input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, group.ResultAdmin{})
}

func GetByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, group.ResultAdmin{
		Group: result.Group[0],
	})
}

func GetAllByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbGroup.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, group.ResultAdminAll{Group: result.Group})
	}
}

func CancelSubscription(c *gin.Context) {
	// ID取得
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// serviceIDが0の時エラー処理
	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("This id is wrong... ")})
		return
	}

	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	resultGroup := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: uint(id)}})
	if resultGroup.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultGroup.Err.Error()})
		return
	}

	if resultGroup.Group[0].StripeSubscriptionID == nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: "Subscription ID is not exists."})
		return
	}

	noticeCancelSubscriptionByAdmin(resultGroup.Group[0])

	stripe.Key = config.Conf.Stripe.SecretKey

	_, err = sub.Cancel(*resultGroup.Group[0].StripeSubscriptionID, nil)
	if err != nil {
		log.Printf("pi.New: %v", err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	err = dbGroup.Update(group.UpdateAll, core.Group{
		StripeSubscriptionID:        &[]string{""}[0],
		PaymentMembershipTemplateID: &[]uint{0}[0],
	})
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.Result{})
}
