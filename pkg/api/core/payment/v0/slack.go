package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/slack-go/slack"
	"github.com/stripe/stripe-go/v72"
	"strconv"
	"strings"
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

func noticePayment(baseKeyValue []string) {
	var slackAttachField []slack.AttachmentField
	for _, keyValue := range baseKeyValue {
		splitKeyValue := strings.Split(keyValue, ":")
		slackAttachField = append(slackAttachField, slack.AttachmentField{Title: splitKeyValue[0], Value: keyValue[len(splitKeyValue[0])+1:]})
	}

	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Payment, slack.MsgOptionAttachments(
		slack.Attachment{
			Color:  "good",
			Title:  "支払い処理",
			Fields: slackAttachField,
		},
	))
}
