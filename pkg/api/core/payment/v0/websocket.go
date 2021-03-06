package v0

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/payment"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbPayment "github.com/homenoc/dsbd-backend/pkg/api/store/payment/v0"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func GetStripeWebHook(c *gin.Context) {
	stripe.Key = config.Conf.Stripe.SecretKey

	const MaxBodyBytes = int64(65536)
	body := http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		return
	}

	event := stripe.Event{}
	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Webhook error while parsing basic request. %v\n", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	endpointSecret := config.Conf.Stripe.WebhookSecretKey

	event, err = webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		log.Println("error: " + err.Error())
		return
	}

	log.Println(event.Type)

	if event.Type == "checkout.session.completed" {
		log.Println("user", event.Data.Object["metadata"].(map[string]interface{})["user"].(string))
	} else if event.Type == "customer.subscription.updated" {
		log.Println("customer.subscription.updated: " + event.Data.Object["id"].(string))
	} else if event.Type == "invoice.paid" {
		log.Println("invoice.paid: " + event.Data.Object["id"].(string))
	} else if event.Type == "invoice.updated" {
		log.Println("invoice.updated: " + event.Data.Object["id"].(string))
	} else if event.Type == "invoice.payment_succeeded" {
		log.Println("invoice.payment_succeeded: " + event.Data.Object["id"].(string))
	} else if event.Type == "payment_intent.created" {
		log.Println("payment_intent.created: " + event.Data.Object["id"].(string))
	} else if event.Type == "payment_intent.succeeded" {
		log.Println("payment_intent.successed: " + event.Data.Object["id"].(string))
		resultPayment, err := dbPayment.Get(payment.PaymentIntentID, core.Payment{PaymentIntentID: event.Data.Object["id"].(string)})
		if err != nil {
			log.Println(err)
		}

		if resultPayment[0].Group == nil {
			log.Println("[paid] error: GroupID is not exists....")
			return
		}

		// Membership
		if *resultPayment[0].IsMembership {
			resultGroup := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: *resultPayment[0].GroupID}})
			if resultGroup.Err != nil {
				log.Println(resultGroup.Err)
				return
			}

			// now time
			now := time.Now()

			if resultGroup.Group[0].MemberExpired != nil {
				now = *resultGroup.Group[0].MemberExpired
			}

			if resultGroup.Group[0].PaymentMembershipTemplate.Yearly {
				// membership yearly
				now = now.AddDate(1, 0, 0)
			} else if resultGroup.Group[0].PaymentMembershipTemplate.Monthly {
				// membership monthly
				now = now.AddDate(0, 1, 0)
			} else {
				log.Println("error:")
				return
			}

			dbGroup.Update(group.UpdateMembership, core.Group{
				Model:         gorm.Model{ID: *resultPayment[0].GroupID},
				MemberExpired: &now,
			})
		}

		err = dbPayment.Update(payment.UpdatePaid, &core.Payment{
			PaymentIntentID: event.Data.Object["id"].(string),
			Paid:            &[]bool{true}[0],
		})
		if err != nil {
			log.Println(err)
		}

		noticeSlackPaymentPaid(event.Data.Object["id"].(string))

	} else if event.Type == "charge.succeeded" {
		log.Printf("charge.succeeded: " + event.Data.Object["id"].(string))
	} else if event.Type == "charge.refunded" {
		log.Printf("charge.succeeded: " + event.Data.Object["id"].(string) + "| " + event.Data.Object["payment_intent"].(string))
		result, err := dbPayment.Get(payment.PaymentIntentID, core.Payment{PaymentIntentID: event.Data.Object["payment_intent"].(string)})
		if err != nil {
			log.Println(err)
		}
		err = dbPayment.Update(payment.UpdateAll, &core.Payment{
			Model:  gorm.Model{ID: result[0].ID},
			Refund: &[]bool{true}[0],
		})
		if err != nil {
			log.Println(err)
		}

	}

	c.JSON(http.StatusOK, common.Result{})
}
