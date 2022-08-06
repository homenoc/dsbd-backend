package v0

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbPayment "github.com/homenoc/dsbd-backend/pkg/api/store/payment/v0"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

	// slack notify(payment log)
	noticePaymentLog(event)

	//t := time.Unix(event.Created, 0)
	//fmt.Println(t)

	switch event.Type {
	case "checkout.session.completed":
		// meta
		_, ok := event.Data.Object["metadata"].(map[string]interface{})["type"]
		if !ok {
			break
		}
		dataType := event.Data.Object["metadata"].(map[string]interface{})["type"].(string)
		name := event.Data.Object["metadata"].(map[string]interface{})["name"].(string)
		groupID := event.Data.Object["metadata"].(map[string]interface{})["group_id"].(string)
		etc := "GroupID: " + groupID + ",  UserName: " + name

		// stripe standard data
		amountTotal := event.Data.Object["amount_total"].(float64)
		paymentIntent := event.Data.Object["payment_intent"].(string)
		if dataType == "donate" {
			dbPayment.Create(&core.Payment{
				Type:            core.PaymentMembership,
				GroupID:         nil,
				Refund:          &[]bool{false}[0],
				PaymentIntentID: paymentIntent,
				Fee:             uint(amountTotal),
			})
			etc += "UserName: " + name
		} else if dataType == "membership" {
			etc += "GroupID: " + groupID
			break
		}

		// slack notify(payment log)
		field := []string{
			"Type:" + dataType + "(" + event.Type + ")",
			"ID:" + event.ID,
			"PaymentIntent:" + paymentIntent,
			"Etc:" + etc,
			"Fee:" + strconv.Itoa(int(uint(amountTotal))) + "円",
		}
		noticePayment(field)
	case "customer.subscription.created":
		// meta
		var dataType, etc string
		_, ok := event.Data.Object["metadata"].(map[string]interface{})["type"]
		if !ok {
			break
		} else {
			dataType = event.Data.Object["metadata"].(map[string]interface{})["type"].(string)
			name := event.Data.Object["metadata"].(map[string]interface{})["name"].(string)
			groupID := event.Data.Object["metadata"].(map[string]interface{})["group_id"].(string)
			etc = "GroupID: " + groupID + ",  UserName: " + name
		}

		// stripe standard data
		customer := event.Data.Object["customer"].(string)
		planID := event.Data.Object["plan"].(map[string]interface{})["id"].(string)
		amount := event.Data.Object["plan"].(map[string]interface{})["amount"].(float64)
		interval := event.Data.Object["plan"].(map[string]interface{})["interval"].(string)
		periodStart := event.Data.Object["current_period_start"].(float64)
		periodEnd := event.Data.Object["current_period_end"].(float64)
		periodStartTime := time.Unix(int64(periodStart), 0)
		periodEndTime := time.Unix(int64(periodEnd), 0)

		// slack notify(payment log)
		field := []string{
			"Type:" + dataType + "(" + event.Type + ")",
			"ID:" + event.ID,
			"CustomerID:" + customer,
			"PlanID:" + planID,
			"Start-EndDate:" + fmt.Sprintf(periodStartTime.Format("2006-01-02")+" - "+periodEndTime.Format("2006-01-02")),
			"Etc:" + etc,
			"Fee:" + strconv.Itoa(int(uint(amount))) + " 円 (" + interval + ")",
		}
		noticePayment(field)
	case "customer.subscription.updated":
		// meta
		dataType := event.Data.Object["metadata"].(map[string]interface{})["type"].(string)
		name := event.Data.Object["metadata"].(map[string]interface{})["name"].(string)
		groupID := event.Data.Object["metadata"].(map[string]interface{})["group_id"].(string)
		etc := "GroupID: " + groupID + ",  UserName: " + name

		// stripe standard data
		customer := event.Data.Object["customer"].(string)
		planID := event.Data.Object["plan"].(map[string]interface{})["id"].(string)
		amount := event.Data.Object["plan"].(map[string]interface{})["amount"].(float64)
		interval := event.Data.Object["plan"].(map[string]interface{})["interval"].(string)
		periodStart := event.Data.Object["current_period_start"].(float64)
		periodEnd := event.Data.Object["current_period_end"].(float64)
		periodStartTime := time.Unix(int64(periodStart), 0)
		periodEndTime := time.Unix(int64(periodEnd), 0)
		status := event.Data.Object["status"].(string)

		// slack notify(payment log)
		field := []string{
			"Type:" + dataType + "(" + event.Type + ")",
			"ID:" + event.ID,
			"CustomerID:" + customer,
			"PlanID:" + planID,
			"Start-EndDate:" + fmt.Sprintf(periodStartTime.Format("2006-01-02")+" - "+periodEndTime.Format("2006-01-02")),
			"Status:" + status,
			"Etc:" + etc,
			"Fee:" + strconv.Itoa(int(uint(amount))) + " 円 (" + interval + ")",
		}
		noticePayment(field)
	case "customer.subscription.deleted":
		// meta
		dataType := event.Data.Object["metadata"].(map[string]interface{})["type"].(string)
		name := event.Data.Object["metadata"].(map[string]interface{})["name"].(string)
		groupID := event.Data.Object["metadata"].(map[string]interface{})["group_id"].(string)
		etc := "GroupID: " + groupID + ",  UserName: " + name

		// stripe standard data
		customer := event.Data.Object["customer"].(string)
		planID := event.Data.Object["plan"].(map[string]interface{})["id"].(string)
		amount := event.Data.Object["plan"].(map[string]interface{})["amount"].(float64)
		interval := event.Data.Object["plan"].(map[string]interface{})["interval"].(string)
		periodStart := event.Data.Object["current_period_start"].(float64)
		periodEnd := event.Data.Object["current_period_end"].(float64)
		periodStartTime := time.Unix(int64(periodStart), 0)
		periodEndTime := time.Unix(int64(periodEnd), 0)
		status := event.Data.Object["status"].(string)

		// slack notify
		field := []string{
			"Type:" + dataType + "(" + event.Type + ")",
			"ID:" + event.ID,
			"CustomerID:" + customer,
			"PlanID:" + planID,
			"Start-EndDate:" + fmt.Sprintf(periodStartTime.Format("2006-01-02")+" - "+periodEndTime.Format("2006-01-02")),
			"Status:" + status,
			"Etc:" + etc,
			"Fee:" + strconv.Itoa(int(uint(amount))) + " 円 (" + interval + ")",
		}
		noticePayment(field)
	}

	c.JSON(http.StatusOK, common.Result{})
}
