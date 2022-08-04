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
	dbPayment "github.com/homenoc/dsbd-backend/pkg/api/store/payment/v0"
	dbPaymentMembershipTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/payment_membership/v0"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/setupintent"
	"github.com/stripe/stripe-go/v72/sub"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func MembershipPayment(c *gin.Context) {
	var input payment.Input
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := auth.GroupAuthorization(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	if result.User.Group.StripeSubscriptionID != nil && *result.User.Group.StripeSubscriptionID != "" {
		c.JSON(http.StatusBadRequest, common.Error{Error: "Error: Subscription."})
		return
	}

	if *result.User.Group.Student {
		c.JSON(http.StatusBadRequest, common.Error{Error: "Error: You are student."})
		return
	}

	resultTemplate, err := dbPaymentMembershipTemplate.Get(input.ItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: "template is not found..."})
		return
	}

	if result.User.Group.StripeCustomerID == nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "stripe customerID is not exists."})
		return
	}

	stripe.Key = config.Conf.Stripe.SecretKey

	params := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(resultTemplate.PriceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
		Customer:        stripe.String(*result.User.Group.StripeCustomerID),
	}
	params.AddExpand("latest_invoice.payment_intent")

	pi, err := sub.New(params)
	log.Printf("pi.New: %v\n", pi.LatestInvoice.PaymentIntent.ClientSecret)
	if err != nil {
		log.Printf("pi.New: %v", err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: "payment_membership system error"})
		return
	}

	dbPayment.Create(&core.Payment{
		GroupID:         result.User.GroupID,
		PaymentIntentID: pi.LatestInvoice.PaymentIntent.ID,
		Type:            core.PaymentMembership,
		Paid:            &[]bool{false}[0],
		Fee:             resultTemplate.Fee,
	})

	dbGroup.Update(group.UpdateMembership, core.Group{
		Model:                       gorm.Model{ID: *result.User.GroupID},
		StripeCustomerID:            result.User.Group.StripeCustomerID,
		StripeSubscriptionID:        &pi.ID,
		PaymentMembershipTemplateID: &resultTemplate.ID,
	})

	go noticeSlackPaymentMembershipPayment(*result.User.GroupID, resultTemplate.Plan, pi.LatestInvoice.PaymentIntent.ID)

	c.JSON(http.StatusOK, payment.ResultByUser{
		ClientSecret: pi.LatestInvoice.PaymentIntent.ClientSecret,
	})
}

func ChangeCardPayment(c *gin.Context) {
	var input payment.ChangeCardPaymentInit
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := auth.GroupAuthorization(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	stripe.Key = config.Conf.Stripe.SecretKey

	if result.User.Group.StripeCustomerID == nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "stripe customerID is not exists."})
		return
	}

	if result.User.Group.StripeSubscriptionID == nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "stripe subscriptionID is not exists."})
		return
	}

	//attach
	pm, err := paymentmethod.Attach(input.PaymentMethodID, &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(*result.User.Group.StripeCustomerID),
	})
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: "payment_membership system error"})
		return
	}
	log.Println(pm)

	// change card on user
	cus, err := customer.Update(*result.User.Group.StripeCustomerID, &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(input.PaymentMethodID),
		},
	})
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: "payment_membership system error"})
		return
	}

	log.Printf(cus.ID)

	_, err = sub.Update(*result.User.Group.StripeSubscriptionID, &stripe.SubscriptionParams{
		DefaultPaymentMethod: stripe.String(input.PaymentMethodID),
	})
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: "payment_membership system error"})
		return
	}
	go noticeSlackPaymentMembershipChangeCardPayment(result.User.Group.ID)

	c.JSON(http.StatusOK, common.Result{})
}

func ChangeCardPaymentInit(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.UserAuthorization(core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	stripe.Key = config.Conf.Stripe.SecretKey

	// get subscription
	params := &stripe.SetupIntentParams{
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
	}

	si, err := setupintent.New(params)
	log.Printf("si.New: %v\n", si.ClientSecret)
	if err != nil {
		log.Printf("si.New: %v", err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: "payment_membership system error"})
		return
	}

	c.JSON(http.StatusOK, payment.ResultByUser{
		ClientSecret: si.ClientSecret,
	})
}
