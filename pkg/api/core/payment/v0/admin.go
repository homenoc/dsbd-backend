package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/payment"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	"github.com/stripe/stripe-go/v73"
	billingSession "github.com/stripe/stripe-go/v73/billingportal/session"
	"github.com/stripe/stripe-go/v73/checkout/session"
	"github.com/stripe/stripe-go/v73/customer"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

func PostAdminSubscribeGettingURL(c *gin.Context) {
	// ID取得
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// Admin authentication
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	// serviceIDが0の時エラー処理
	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("ID is wrong... ")})
		return
	}

	result := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	var input payment.Input
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// search plan
	membershipWithTemplate, err := config.GetMembershipTemplate(input.Plan)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: "invalid plan"})
		return
	}

	// exist check: stripeCustomerID
	if result.Group[0].StripeCustomerID == nil || *result.Group[0].StripeCustomerID == "" {
		params := &stripe.CustomerParams{
			Description: stripe.String("[" + strconv.Itoa(int(result.Group[0].ID)) + "] Org: " + result.Group[0].Org + "(" + result.Group[0].OrgEn + ")"),
		}
		cus, err := customer.New(params)
		if err != nil {
			noticePaymentError(true, []string{"Type: Create Customer", "Error: " + err.Error()})
			log.Println("Error: " + err.Error())
		}
		err = dbGroup.Update(group.UpdateAll, core.Group{Model: gorm.Model{ID: result.Group[0].ID}, StripeCustomerID: &cus.ID})
		noticePaymentLog(stripe.Event{
			ID:   cus.ID,
			Type: "stripe customer追加",
		})
		result.Group[0].StripeCustomerID = &cus.ID
	}

	date := time.Now()
	params := &stripe.CheckoutSessionParams{
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Customer: result.Group[0].StripeCustomerID,
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(membershipWithTemplate.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(config.Conf.Controller.Admin.ReturnURL),
		CancelURL:  stripe.String(config.Conf.Controller.Admin.ReturnURL),
		ExpiresAt:  stripe.Int64(date.Add(time.Minute * 30).Unix()),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"type":     "membership",
				"name":     "---",
				"group_id": strconv.Itoa(int(result.Group[0].ID)),
				"log": "[---] admin" +
					"_[" + strconv.Itoa(int(result.Group[0].ID)) + "] " + result.Group[0].Org,
			},
		},
	}

	s, err := session.New(params)

	if err != nil {
		noticePaymentError(true, []string{"Type: Subscribe Session", "Error: " + err.Error()})
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"url": s.URL})
}

func GetAdminBillingPortalURL(c *gin.Context) {
	// ID取得
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// Admin authentication
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	// serviceIDが0の時エラー処理
	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("ID is wrong... ")})
		return
	}

	result := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	// exist check: stripeCustomerID
	if result.Group[0].StripeCustomerID == nil || *result.Group[0].StripeCustomerID == "" {
		c.JSON(http.StatusNotFound, common.Error{Error: "CustomerID is not found..."})
		return
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:      stripe.String(*result.Group[0].StripeCustomerID),
		Configuration: stripe.String(config.Conf.Stripe.MembershipConfiguration),
		ReturnURL:     stripe.String(config.Conf.Controller.Admin.ReturnURL),
	}

	s, err := billingSession.New(params)
	if err != nil {
		noticePaymentError(true, []string{"Type: BillingSession", "Error: " + err.Error()})
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"url": s.URL})
}
