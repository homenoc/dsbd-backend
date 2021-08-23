package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"gorm.io/gorm"
	"strconv"
)

func noticeSlackPaymentMembershipPayment(groupID uint, plan, paymentIntentID string) {
	attachment := slack.Attachment{}

	result := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: groupID}})

	attachment.Text = &[]string{"会費支払い"}[0]
	attachment.AddField(slack.Field{Title: "申請者", Value: strconv.Itoa(int(groupID)) + ": " + result.Group[0].Org + "(" + result.Group[0].OrgEn + ")"}).
		AddField(slack.Field{Title: "Plan", Value: plan}).
		AddField(slack.Field{Title: "PaymentIntentID", Value: paymentIntentID})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackPaymentDonatePayment(userID, money uint, paymentIntentID string) {
	attachment := slack.Attachment{}

	result := dbUser.Get(user.ID, &core.User{Model: gorm.Model{ID: userID}})

	attachment.Text = &[]string{"寄付"}[0]
	attachment.AddField(slack.Field{Title: "申請者", Value: strconv.Itoa(int(userID)) + ": " + result.User[0].Name + "(" + result.User[0].NameEn + ")"}).
		AddField(slack.Field{Title: "金額", Value: strconv.Itoa(int(money)) + "円"}).
		AddField(slack.Field{Title: "PaymentIntentID", Value: paymentIntentID})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackPaymentMembershipChangeCardPayment(groupID uint) {
	attachment := slack.Attachment{}

	result := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: groupID}})

	attachment.Text = &[]string{"会費支払い(カードの変更)"}[0]
	attachment.AddField(slack.Field{Title: "申請者", Value: strconv.Itoa(int(groupID)) + ": " + result.Group[0].Org + "(" + result.Group[0].OrgEn + ")"})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackPaymentPaid(paymentIntentID string) {
	attachment := slack.Attachment{}

	attachment.Text = &[]string{"支払い完了"}[0]
	attachment.AddField(slack.Field{Title: "PaymentIntentID", Value: paymentIntentID})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}
