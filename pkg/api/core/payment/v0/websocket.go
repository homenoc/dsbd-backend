package v0

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/webhook"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetStripeWebHook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	body := http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		return
	}

	event := stripe.Event{}
	if err := json.Unmarshal(payload, &event); err != nil {
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
		groupIDStr := event.Data.Object["metadata"].(map[string]interface{})["group_id"].(string)
		etc := "GroupID: " + groupIDStr + ",  UserName: " + name

		// stripe standard data
		amountTotal := event.Data.Object["amount_total"].(float64)
		paymentIntent := event.Data.Object["payment_intent"].(string)
		if dataType == "donate" {
			etc += "UserName: " + name
		} else if dataType == "donate_membership" {
			etc += "UserName: " + name
		} else if dataType == "membership" {
			etc += "GroupID: " + groupIDStr
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
		var dataType, groupIDStr, etc string
		var groupID int
		_, ok := event.Data.Object["metadata"].(map[string]interface{})["type"]
		if !ok {
			break
		} else {
			dataType = event.Data.Object["metadata"].(map[string]interface{})["type"].(string)
			name := event.Data.Object["metadata"].(map[string]interface{})["name"].(string)
			groupIDStr = event.Data.Object["metadata"].(map[string]interface{})["group_id"].(string)
			etc = "GroupID: " + groupIDStr + ",  UserName: " + name
		}

		if dataType == "membership" {
			groupID, _ = strconv.Atoi(groupIDStr)

			// stripe standard data
			customer := event.Data.Object["customer"].(string)
			sub := event.Data.Object["id"].(string)
			planID := event.Data.Object["plan"].(map[string]interface{})["id"].(string)
			amount := event.Data.Object["plan"].(map[string]interface{})["amount"].(float64)
			interval := event.Data.Object["plan"].(map[string]interface{})["interval"].(string)
			periodStart := event.Data.Object["current_period_start"].(float64)
			periodEnd := event.Data.Object["current_period_end"].(float64)
			periodStartTime := time.Unix(int64(periodStart), 0)
			periodEndTime := time.Unix(int64(periodEnd), 0)
			jst, _ := time.LoadLocation(config.Conf.Controller.TimeZone)
			timeDate := time.Date(periodEndTime.Year(), periodEndTime.Month(), periodEndTime.Day(), 0, 0, 0, 0, jst)
			if groupID != 0 {
				resultGroup := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: uint(groupID)}})
				if resultGroup.Err != nil {
					return
				}
				if resultGroup.Group[0].MemberExpired.Unix() < timeDate.Unix() {
					err = dbGroup.Update(group.UpdateAll, core.Group{Model: gorm.Model{ID: uint(groupID)}, StripeSubscriptionID: &sub, MemberExpired: &timeDate})
				}
			}
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
		} else if dataType == "donate_membership" {
			customer := event.Data.Object["customer"].(string)
			planID := event.Data.Object["plan"].(map[string]interface{})["id"].(string)
			amount := event.Data.Object["plan"].(map[string]interface{})["amount"].(float64)
			periodStart := event.Data.Object["current_period_start"].(float64)
			periodEnd := event.Data.Object["current_period_end"].(float64)
			periodStartTime := time.Unix(int64(periodStart), 0)
			periodEndTime := time.Unix(int64(periodEnd), 0)

			// slack notify(payment log)
			field := []string{
				"Type: [寄付]" + dataType + "(" + event.Type + ")",
				"ID:" + event.ID,
				"CustomerID:" + customer,
				"PlanID:" + planID,
				"Start-EndDate:" + fmt.Sprintf(periodStartTime.Format("2006-01-02")+" - "+periodEndTime.Format("2006-01-02")),
				"Etc:" + etc,
				"Fee:" + strconv.Itoa(int(uint(amount))) + " 円",
			}
			noticePayment(field)
		}
	case "customer.subscription.updated":
		var groupID int

		// meta
		dataType := event.Data.Object["metadata"].(map[string]interface{})["type"].(string)
		name := event.Data.Object["metadata"].(map[string]interface{})["name"].(string)
		groupIDStr := event.Data.Object["metadata"].(map[string]interface{})["group_id"].(string)
		etc := "GroupID: " + groupIDStr + ",  UserName: " + name
		if dataType == "membership" {
			groupID, _ = strconv.Atoi(groupIDStr)

			// stripe standard data
			sub := event.Data.Object["id"].(string)
			customer := event.Data.Object["customer"].(string)
			planID := event.Data.Object["plan"].(map[string]interface{})["id"].(string)
			amount := event.Data.Object["plan"].(map[string]interface{})["amount"].(float64)
			interval := event.Data.Object["plan"].(map[string]interface{})["interval"].(string)
			periodStart := event.Data.Object["current_period_start"].(float64)
			periodEnd := event.Data.Object["current_period_end"].(float64)
			periodStartTime := time.Unix(int64(periodStart), 0)
			periodEndTime := time.Unix(int64(periodEnd), 0)
			status := event.Data.Object["status"].(string)
			jst, _ := time.LoadLocation(config.Conf.Controller.TimeZone)
			timeDate := time.Date(periodEndTime.Year(), periodEndTime.Month(), periodEndTime.Day(), 0, 0, 0, 0, jst)
			if groupID != 0 {
				err = dbGroup.Update(group.UpdateAll, core.Group{Model: gorm.Model{ID: uint(groupID)}, StripeSubscriptionID: &sub, MemberExpired: &timeDate})
			}

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
		}
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
