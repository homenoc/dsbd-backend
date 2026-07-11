package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/payment"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	"github.com/stripe/stripe-go/v73"
	billingSession "github.com/stripe/stripe-go/v73/billingportal/session"
	"github.com/stripe/stripe-go/v73/checkout/session"
	"github.com/stripe/stripe-go/v73/customer"
	"log"
	"net/http"
	"strconv"
	"time"
)

func PostSubscribeGettingURL(c *gin.Context) {
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

	user := middleware.CurrentUser(c)

	// exist check: stripeCustomerID
	if user.Group.StripeCustomerID == nil || *user.Group.StripeCustomerID == "" {
		params := &stripe.CustomerParams{
			Description: stripe.String("[" + strconv.Itoa(int(user.Group.ID)) + "] Org: " + user.Group.Org + "(" + user.Group.OrgEn + ")"),
		}
		cus, err := customer.New(params)
		if err != nil {
			noticePaymentError(false, []string{
				"User: [" + strconv.Itoa(int(user.ID)) + "] " + user.Name,
				"Group: [" + strconv.Itoa(int(user.Group.ID)) + "] " + user.Group.Org,
				"Type: Create Customer", "Error: " + err.Error()},
			)
			log.Println("Error: " + err.Error())
		}
		err = dbGroup.UpdateStripeCustomerID(user.Group.ID, cus.ID)
		noticePaymentLog(stripe.Event{
			ID:   cus.ID,
			Type: "stripe customer追加",
		})
		user.Group.StripeCustomerID = &cus.ID
	}

	date := time.Now()
	params := &stripe.CheckoutSessionParams{
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Customer: user.Group.StripeCustomerID,
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
				"group_id": strconv.Itoa(int(user.Group.ID)),
				"name":     strconv.Itoa(int(user.ID)),
				"log": "[" + strconv.Itoa(int(user.ID)) + "] " + user.Name +
					"_[" + strconv.Itoa(int(user.Group.ID)) + "] " + user.Group.Org,
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
	user := middleware.CurrentUser(c)

	// exist check: stripeCustomerID
	if user.Group.StripeCustomerID == nil || *user.Group.StripeCustomerID == "" {
		c.JSON(http.StatusNotFound, common.Error{Error: "CustomerID is not found..."})
		return
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(*user.Group.StripeCustomerID),
		ReturnURL: stripe.String(config.Conf.Controller.User.ReturnURL),
	}

	s, err := billingSession.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"url": s.URL})
}
