package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/payment"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/stripe/stripe-go/v73"
	billingSession "github.com/stripe/stripe-go/v73/billingportal/session"
	"github.com/stripe/stripe-go/v73/checkout/session"
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
	if *resultAuth.User.Group.StripeCustomerID == "" {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "stripe customer id not found..."})
		return
	}

	// exist check: stripeCustomerID
	if *resultAuth.User.Group.StripeCustomerID == "" {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "stripe customer id not found..."})
		return
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
		SuccessURL: stripe.String(config.Conf.Controller.User.Url),
		//CancelURL:  stripe.String("https://example.com/cancel"),
		ExpiresAt: stripe.Int64(date.Add(time.Minute * 30).Unix()),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"type":     "membership",
				"group_id": strconv.Itoa(int(resultAuth.User.Group.ID)),
				"name":     "Yuto Yoneda",
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

	params := &stripe.BillingPortalSessionParams{
		Customer:      stripe.String("cus_MBBRylUHxUQvVc"),
		Configuration: stripe.String(config.Conf.Stripe.MembershipConfiguration),
		ReturnURL:     stripe.String(config.Conf.Controller.User.Url),
	}

	s, err := billingSession.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"url": s.URL})
}
