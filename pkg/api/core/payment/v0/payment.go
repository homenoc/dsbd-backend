package v0

import (
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

func PostSubscribeGettingURL(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	var input payment.Input
	err := c.BindJSON(&input)
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

	resultAuth := auth.GroupAuthorization(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if resultAuth.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAuth.Err.Error()})
		return
	}

	// exist check: stripeCustomerID
	if resultAuth.User.Group.StripeCustomerID == nil || *resultAuth.User.Group.StripeCustomerID == "" {
		params := &stripe.CustomerParams{
			Description: stripe.String("[" + strconv.Itoa(int(resultAuth.User.Group.ID)) + "] Org: " + resultAuth.User.Group.Org + "(" + resultAuth.User.Group.OrgEn + ")"),
		}
		cus, err := customer.New(params)
		if err != nil {
			noticePaymentError(false, []string{
				"User: [" + strconv.Itoa(int(resultAuth.User.ID)) + "] " + resultAuth.User.Name,
				"Group: [" + strconv.Itoa(int(resultAuth.User.Group.ID)) + "] " + resultAuth.User.Group.Org,
				"Type: Create Customer", "Error: " + err.Error()},
			)
			log.Println("Error: " + err.Error())
		}
		err = dbGroup.Update(group.UpdateAll, core.Group{Model: gorm.Model{ID: resultAuth.User.Group.ID}, StripeCustomerID: &cus.ID})
		noticePaymentLog(stripe.Event{
			ID:   cus.ID,
			Type: "stripe customer追加",
		})
		resultAuth.User.Group.StripeCustomerID = &cus.ID
	}

	date := time.Now()
	params := &stripe.CheckoutSessionParams{
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Customer: resultAuth.User.Group.StripeCustomerID,
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(membershipWithTemplate.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(config.Conf.Controller.User.ReturnURL),
		CancelURL:  stripe.String(config.Conf.Controller.User.ReturnURL),
		ExpiresAt:  stripe.Int64(date.Add(time.Minute * 30).Unix()),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"type":     "membership",
				"group_id": strconv.Itoa(int(resultAuth.User.Group.ID)),
				"name":     strconv.Itoa(int(resultAuth.User.ID)),
				"log": "[" + strconv.Itoa(int(resultAuth.User.ID)) + "] " + resultAuth.User.Name +
					"_[" + strconv.Itoa(int(resultAuth.User.Group.ID)) + "] " + resultAuth.User.Group.Org,
			},
		},
	}

	s, err := session.New(params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"url": s.URL})
}

func GetBillingPortalURL(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	resultAuth := auth.GroupAuthorization(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if resultAuth.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAuth.Err.Error()})
		return
	}

	// exist check: stripeCustomerID
	if resultAuth.User.Group.StripeCustomerID == nil || *resultAuth.User.Group.StripeCustomerID == "" {
		c.JSON(http.StatusNotFound, common.Error{Error: "CustomerID is not found..."})
		return
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(*resultAuth.User.Group.StripeCustomerID),
		ReturnURL: stripe.String(config.Conf.Controller.User.ReturnURL),
	}

	s, err := billingSession.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"url": s.URL})
}
