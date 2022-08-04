package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/mail"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/slack-go/slack"
)

func noticeSlack(err error, input mail.Mail) {
	// 審査ステータスのSlack通知
	if err != nil {
		notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Log, slack.MsgOptionAttachments(
			slack.Attachment{
				Color: "danger",
				Title: "メール送信(失敗)",
				Text:  "error: \n" + err.Error(),
				Fields: []slack.AttachmentField{
					{Title: "Subject", Value: input.Subject},
					{Title: "Content", Value: input.Content},
				},
			},
		))
	} else {
		notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Log, slack.MsgOptionAttachments(
			slack.Attachment{
				Color: "good",
				Title: "メール送信",
				Fields: []slack.AttachmentField{
					{Title: "Subject", Value: input.Subject},
					{Title: "Content", Value: input.Content},
				},
			},
		))
	}
}
