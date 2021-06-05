package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core/mail"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
)

func noticeSlack(err error, input mail.Mail) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	if err != nil {
		attachment.AddField(slack.Field{Title: "Title", Value: "メール送信(失敗)"}).
			AddField(slack.Field{Title: "To", Value: input.ToMail}).
			AddField(slack.Field{Title: "Subject", Value: input.Subject}).
			AddField(slack.Field{Title: "Content", Value: input.Content}).
			AddField(slack.Field{Title: "Error", Value: err.Error()})
		notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: false})
	} else {
		//attachment.AddField(slack.Field{Title: "Title", Value: "メール送信"}).
		//	AddField(slack.Field{Title: "To", Value: input.ToMail}).
		//	AddField(slack.Field{Title: "Subject", Value: input.Subject}).
		//	AddField(slack.Field{Title: "Content", Value: input.Content})
		//notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
	}
}
