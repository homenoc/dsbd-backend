package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/slack-go/slack"
	"github.com/stripe/stripe-go/v72"
	"strconv"
)

func noticePaymentLog(event stripe.Event) {
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.PaymentLog, slack.MsgOptionAttachments(
		slack.Attachment{
			Color: "good",
			Title: event.Type,
			Fields: []slack.AttachmentField{
				{Title: "ID", Value: event.ID},
				{Title: "Created", Value: strconv.FormatInt(event.Created, 10)},
			},
		},
	))
}

func noticePayment(keyValue map[string]string) {
	var slackAttachField []slack.AttachmentField
	for key, value := range keyValue {
		slackAttachField = append(slackAttachField, slack.AttachmentField{Title: key, Value: value})
	}

	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.PaymentLog, slack.MsgOptionAttachments(
		slack.Attachment{
			Color:  "good",
			Title:  "支払い処理",
			Fields: slackAttachField,
		},
	))
}
